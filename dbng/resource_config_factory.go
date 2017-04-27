package dbng

import (
	"database/sql"

	"code.cloudfoundry.org/lager"
	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db/lock"
	"github.com/lib/pq"
)

//go:generate counterfeiter . ResourceConfigFactory

type ResourceConfigFactory interface {
	FindOrCreateResourceConfig(
		logger lager.Logger,
		user ResourceUser,
		resourceType string,
		source atc.Source,
		resourceTypes atc.VersionedResourceTypes,
	) (*UsedResourceConfig, error)

	CleanConfigUsesForFinishedBuilds(lager.Logger) error
	CleanConfigUsesForInactiveResourceTypes(lager.Logger) error
	CleanConfigUsesForInactiveResources(lager.Logger) error
	CleanConfigUsesForPausedPipelinesResources(lager.Logger) error
	CleanConfigUsesForOutdatedResourceConfigs(lager.Logger) error
	CleanUselessConfigs(lager.Logger) error

	AcquireResourceCheckingLock(
		logger lager.Logger,
		resourceUser ResourceUser,
		resourceType string,
		resourceSource atc.Source,
		resourceTypes atc.VersionedResourceTypes,
	) (lock.Lock, bool, error)
}

type resourceConfigFactory struct {
	conn        Conn
	lockFactory lock.LockFactory
}

func NewResourceConfigFactory(conn Conn, lockFactory lock.LockFactory) ResourceConfigFactory {
	return &resourceConfigFactory{
		conn:        conn,
		lockFactory: lockFactory,
	}
}

func (f *resourceConfigFactory) FindOrCreateResourceConfig(
	logger lager.Logger,
	user ResourceUser,
	resourceType string,
	source atc.Source,
	resourceTypes atc.VersionedResourceTypes,
) (*UsedResourceConfig, error) {
	resourceConfig, err := constructResourceConfig(resourceType, source, resourceTypes)
	if err != nil {
		return nil, err
	}

	var usedResourceConfig *UsedResourceConfig

	err = safeFindOrCreate(f.conn, func(tx Tx) error {
		var err error

		usedResourceConfig, err = user.UseResourceConfig(logger, tx, f.lockFactory, resourceConfig)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return usedResourceConfig, nil
}

func (f *resourceConfigFactory) CleanConfigUsesForFinishedBuilds(logger lager.Logger) error {
	return f.logAndDeleteUses(
		logger,
		psql.Delete("resource_config_uses rcu USING builds b").
			Where(sq.Expr("rcu.build_id = b.id")).
			Where(sq.Expr("NOT b.interceptible")),
	)
}

func (f *resourceConfigFactory) CleanConfigUsesForInactiveResourceTypes(logger lager.Logger) error {
	return f.logAndDeleteUses(
		logger,
		psql.Delete("resource_config_uses rcu USING resource_types t").
			Where(sq.And{
				sq.Expr("rcu.resource_type_id = t.id"),
				sq.Eq{
					"t.active": false,
				},
			}),
	)
}

func (f *resourceConfigFactory) CleanConfigUsesForInactiveResources(logger lager.Logger) error {
	return f.logAndDeleteUses(
		logger,
		psql.Delete("resource_config_uses rcu USING resources r").
			Where(sq.And{
				sq.Expr("rcu.resource_id = r.id"),
				sq.Eq{
					"r.active": false,
				},
			}),
	)
}

func (f *resourceConfigFactory) CleanConfigUsesForPausedPipelinesResources(logger lager.Logger) error {
	pausedPipelineIds, _, err := sq.
		Select("id").
		Distinct().
		From("pipelines").
		Where(sq.Expr("paused = false")).
		ToSql()
	if err != nil {
		return err
	}

	return f.logAndDeleteUses(
		logger,
		psql.Delete("resource_config_uses rcu USING resources r").
			Where(sq.And{
				sq.Expr("r.pipeline_id NOT IN (" + pausedPipelineIds + ")"),
				sq.Expr("rcu.resource_id = r.id"),
			}),
	)
}

func (f *resourceConfigFactory) CleanConfigUsesForOutdatedResourceConfigs(logger lager.Logger) error {
	return f.logAndDeleteUses(
		logger,
		psql.Delete("resource_config_uses rcu USING resources r, resource_configs rc").
			Where(sq.And{
				sq.Expr("rcu.resource_id = r.id"),
				sq.Expr("rcu.resource_config_id = rc.id"),
				sq.Expr("r.source_hash != rc.source_hash"),
			}),
	)
}

func (f *resourceConfigFactory) CleanUselessConfigs(logger lager.Logger) error {
	stillInUseConfigIds, _, err := sq.
		Select("resource_config_id").
		Distinct().
		From("resource_config_uses").
		ToSql()
	if err != nil {
		return err
	}

	usedByResourceCachesIds, _, err := sq.
		Select("resource_config_id").
		Distinct().
		From("resource_caches").
		ToSql()
	if err != nil {
		return err
	}

	delete := psql.Delete("resource_configs").
		Where("id NOT IN (" + stillInUseConfigIds + ")").
		Where("id NOT IN (" + usedByResourceCachesIds + ")").
		Suffix("RETURNING id, base_resource_type_id, resource_cache_id, source_hash").
		PlaceholderFormat(sq.Dollar)

	rows, err := sq.QueryWith(f.conn, delete)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "foreign_key_violation" {
			// this can happen if a use or resource cache is created referencing the
			// config; as the subqueries above are not atomic
			return nil
		}

		return err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var baseResourceTypeID, resourceCacheID sql.NullInt64
		var sourceHash string
		err := rows.Scan(&id, &baseResourceTypeID, &resourceCacheID, &sourceHash)
		if err != nil {
			logger.Error("failed-to-scan-deleted-row", err)
			return err
		}

		data := lager.Data{
			"id":          id,
			"source-hash": sourceHash,
		}

		if baseResourceTypeID.Valid {
			data["base-resource-type-id"] = baseResourceTypeID.Int64
		}

		if resourceCacheID.Valid {
			data["resource-cache-id"] = resourceCacheID.Int64
		}

		logger.Debug("deleted-resource-config", data)
	}

	return nil
}

func resourceTypesList(resourceTypeName string, allResourceTypes []atc.ResourceType, resultResourceTypes []atc.ResourceType) []atc.ResourceType {
	for _, resourceType := range allResourceTypes {
		if resourceType.Name == resourceTypeName {
			resultResourceTypes = append(resultResourceTypes, resourceType)
			return resourceTypesList(resourceType.Type, allResourceTypes, resultResourceTypes)
		}
	}

	return resultResourceTypes
}

func (f *resourceConfigFactory) AcquireResourceCheckingLock(
	logger lager.Logger,
	resourceUser ResourceUser,
	resourceType string,
	resourceSource atc.Source,
	resourceTypes atc.VersionedResourceTypes,
) (lock.Lock, bool, error) {
	resourceConfig, err := constructResourceConfig(resourceType, resourceSource, resourceTypes)
	if err != nil {
		return nil, false, err
	}

	logger.Debug("acquiring-resource-checking-lock", lager.Data{
		"resource-config": resourceConfig,
		"resource-type":   resourceType,
		"resource-source": resourceSource,
		"resource-types":  resourceTypes,
	})

	return acquireResourceCheckingLock(
		logger.Session("lock", lager.Data{"resource-user": resourceUser}),
		f.conn,
		resourceUser,
		resourceConfig,
		f.lockFactory,
	)
}

func (f *resourceConfigFactory) logAndDeleteUses(logger lager.Logger, delete sq.DeleteBuilder) error {
	delete = delete.Suffix("RETURNING resource_config_id, build_id, resource_id, resource_type_id")

	rows, err := sq.QueryWith(f.conn, delete)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var resourceConfigID int
		var buildID, resourceID, resourceTypeID sql.NullInt64
		err := rows.Scan(&resourceConfigID, &buildID, &resourceID, &resourceTypeID)
		if err != nil {
			logger.Error("failed-to-scan-deleted-row", err)
			return err
		}

		data := lager.Data{
			"resource-config-id": resourceConfigID,
		}

		if buildID.Valid {
			data["build-id"] = buildID.Int64
		}

		if resourceID.Valid {
			data["resource-id"] = resourceID.Int64
		}

		if resourceTypeID.Valid {
			data["resource-type-id"] = resourceTypeID.Int64
		}

		logger.Debug("deleted-resource-config-use", data)
	}

	return nil
}

func constructResourceConfig(
	resourceType string,
	source atc.Source,
	resourceTypes atc.VersionedResourceTypes,
) (ResourceConfig, error) {
	resourceConfig := ResourceConfig{
		Source: source,
	}

	customType, found := resourceTypes.Lookup(resourceType)
	if found {
		customTypeResourceConfig, err := constructResourceConfig(
			customType.Type,
			customType.Source,
			resourceTypes.Without(customType.Name),
		)
		if err != nil {
			return ResourceConfig{}, err
		}

		resourceConfig.CreatedByResourceCache = &ResourceCache{
			ResourceConfig: customTypeResourceConfig,
			Version:        customType.Version,
		}
	} else {
		resourceConfig.CreatedByBaseResourceType = &BaseResourceType{
			Name: resourceType,
		}
	}

	return resourceConfig, nil
}

func acquireResourceCheckingLock(
	logger lager.Logger,
	conn Conn,
	user ResourceUser,
	resourceConfig ResourceConfig,
	lockFactory lock.LockFactory,
) (lock.Lock, bool, error) {
	var usedResourceConfig *UsedResourceConfig

	err := safeFindOrCreate(conn, func(tx Tx) error {
		var err error

		usedResourceConfig, err = user.UseResourceConfig(logger, tx, lockFactory, resourceConfig)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	lock := lockFactory.NewLock(
		logger,
		lock.NewResourceConfigCheckingLockID(usedResourceConfig.ID),
	)

	acquired, err := lock.Acquire()
	if err != nil {
		return nil, false, err
	}

	if !acquired {
		return nil, false, nil
	}

	return lock, true, nil
}
