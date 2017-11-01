package radar

import (
	"log"
	"reflect"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

type resourceTypeScanner struct {
	clock                             clock.Clock
	resourceFactory                   resource.ResourceFactory
	resourceConfigCheckSessionFactory db.ResourceConfigCheckSessionFactory
	defaultInterval                   time.Duration
	dbPipeline                        db.Pipeline
	variables                         creds.Variables
}

func NewResourceTypeScanner(
	clock clock.Clock,
	resourceFactory resource.ResourceFactory,
	resourceConfigCheckSessionFactory db.ResourceConfigCheckSessionFactory,
	defaultInterval time.Duration,
	dbPipeline db.Pipeline,
	variables creds.Variables,
) Scanner {
	return &resourceTypeScanner{
		clock:                             clock,
		resourceFactory:                   resourceFactory,
		resourceConfigCheckSessionFactory: resourceConfigCheckSessionFactory,
		defaultInterval:                   defaultInterval,
		dbPipeline:                        dbPipeline,
		variables:                         variables,
	}
}

func (scanner *resourceTypeScanner) Run(logger lager.Logger, resourceTypeName string) (time.Duration, error) {
	return scanner.prescan(logger.Session("tick"), resourceTypeName, nil, false)
}

func (scanner *resourceTypeScanner) ScanFromVersion(logger lager.Logger, resourceTypeName string, fromVersion atc.Version) error {
	// FIXME: Implement
	return nil
}

func (scanner *resourceTypeScanner) Scan(logger lager.Logger, resourceTypeName string) error {
	_, err := scanner.prescan(logger.Session("tick"), resourceTypeName, nil, true)
	return err
}

func (scanner *resourceTypeScanner) prescan(logger lager.Logger, resourceTypeName string, fromVersion atc.Version, mustComplete bool) (time.Duration, error) {
	// load up resource types from db
	// return only the set required by the given type (transiviely, i.e. if its custom type depends on another, include it)
	// evaluates their credentials
	// for any that don't ahve a version, perform a check, and save their version
	// RETURNS atc.VersionedResourceTypes - we can probably remove creds.VersionedResourceTypes
	savedResourceType, found, err := scanner.dbPipeline.ResourceType(resourceTypeName)
	if err != nil {
		logger.Error("failed-to-get-current-version", err)
		return 0, err
	}

	if !found {
		return 0, db.ResourceTypeNotFoundError{Name: resourceTypeName}
	}

	resourceTypesFactory := NewResourceTypeFactory(scanner.dbPipeline, scanner.variables)
	resourceTypeDependencies, err := resourceTypesFactory.ResourceTypes(logger, savedResourceType.Type())
	if err != nil {
		// XXX test
		return 0, err
	}
	for _, rtDependency := range resourceTypeDependencies {
		if rtDependency.Version == nil {
			_, err := scanner.scan(logger, rtDependency.ResourceType.Name, fromVersion, true)
			log.Println("SCAN", rtDependency.ResourceType.Name, err)
			if err != nil {
				logger.Error("failed-to-scan-dependency", err, lager.Data{"depName": rtDependency.ResourceType.Name})
			}
		}
	}

	return scanner.scan(logger, resourceTypeName, fromVersion, mustComplete)
}

func (scanner *resourceTypeScanner) scan(logger lager.Logger, resourceTypeName string, fromVersion atc.Version, mustComplete bool) (time.Duration, error) {
	lockLogger := logger.Session("lock", lager.Data{
		"resource-type": resourceTypeName,
	})

	savedResourceType, found, err := scanner.dbPipeline.ResourceType(resourceTypeName)
	if err != nil {
		logger.Error("failed-to-get-current-version", err)
		return 0, err
	}

	if !found {
		return 0, db.ResourceTypeNotFoundError{Name: resourceTypeName}
	}

	// TODO: maybe consider scanner.checkInterval
	interval := scanner.defaultInterval

	resourceTypes, err := scanner.dbPipeline.ResourceTypes()
	if err != nil {
		logger.Error("failed-to-get-resource-types", err)
		return 0, err
	}

	// // FIXME: Scan dependencies
	// //   Go through each resourceType
	// //     if the resourceType's Name matches the savedResourceType's Type
	// //       if the resoureceType's Version is nil
	// //         Scan the resourceType
	// //   Reload all the resourceTypes at some point

	versionedResourceTypes := creds.NewVersionedResourceTypes(
		scanner.variables,
		resourceTypes.Deserialize(),
	)

	source, err := creds.NewSource(scanner.variables, savedResourceType.Source()).Evaluate()
	if err != nil {
		logger.Error("failed-to-evaluate-resource-type-source", err)
		return 0, err
	}

	resourceConfigCheckSession, err := scanner.resourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSession(
		logger,
		savedResourceType.Type(),
		source,
		versionedResourceTypes.Without(savedResourceType.Name()), //<- may no longer be necessary as the factory knows to exclude given type - check after!
		ContainerExpiries,
	)
	if err != nil {
		logger.Error("failed-to-find-or-create-resource-config", err)
		return 0, err
	}

	err = savedResourceType.SetResourceConfig(resourceConfigCheckSession.ResourceConfig().ID)
	if err != nil {
		logger.Error("failed-to-set-resource-config-id-on-resource-type", err)
		return 0, err
	}

	for breaker := true; breaker == true; breaker = mustComplete {
		lock, acquired, err := scanner.dbPipeline.AcquireResourceTypeCheckingLockWithIntervalCheck(
			logger,
			savedResourceType.Name(),
			resourceConfigCheckSession.ResourceConfig(),
			interval,
			mustComplete,
		)
		if err != nil {
			lockLogger.Error("failed-to-get-lock", err, lager.Data{
				"resource-type": resourceTypeName,
			})
			return interval, ErrFailedToAcquireLock
		}

		if !acquired {
			lockLogger.Debug("did-not-get-lock")
			if mustComplete {
				scanner.clock.Sleep(time.Second)
				continue
			} else {
				return interval, ErrFailedToAcquireLock
			}
		}

		defer lock.Release()

		break
	}

	if fromVersion == nil {
		fromVersion = atc.Version(savedResourceType.Version())
	}

	return interval, scanner.check(
		logger,
		savedResourceType,
		resourceConfigCheckSession,
		fromVersion,
		versionedResourceTypes,
		source,
	)
}

func (scanner *resourceTypeScanner) check(
	logger lager.Logger,
	savedResourceType db.ResourceType,
	resourceConfigCheckSession db.ResourceConfigCheckSession,
	fromVersion atc.Version,
	versionedResourceTypes creds.VersionedResourceTypes,
	source atc.Source,
) error {
	pipelinePaused, err := scanner.dbPipeline.CheckPaused()
	if err != nil {
		logger.Error("failed-to-check-if-pipeline-paused", err)
		return err
	}

	if pipelinePaused {
		logger.Debug("pipeline-paused")
		return nil
	}

	resourceSpec := worker.ContainerSpec{
		ImageSpec: worker.ImageSpec{
			ResourceType: savedResourceType.Type(),
		},
		Tags:   []string{},
		TeamID: scanner.dbPipeline.TeamID(),
	}

	res, err := scanner.resourceFactory.NewResource(
		logger,
		nil,
		db.NewResourceConfigCheckSessionContainerOwner(resourceConfigCheckSession, scanner.dbPipeline.TeamID()),
		db.ContainerMetadata{
			Type: db.ContainerTypeCheck,
		},
		resourceSpec,
		versionedResourceTypes.Without(savedResourceType.Name()),
		worker.NoopImageFetchingDelegate{},
	)
	if err != nil {
		logger.Error("failed-to-initialize-new-container", err)
		return err
	}

	newVersions, err := res.Check(source, fromVersion)
	if err != nil {
		if rErr, ok := err.(resource.ErrResourceScriptFailed); ok {
			logger.Info("check-failed", lager.Data{"exit-status": rErr.ExitStatus})
			return nil
		}

		logger.Error("failed-to-check", err)
		return err
	}

	if len(newVersions) == 0 || reflect.DeepEqual(newVersions, []atc.Version{fromVersion}) {
		logger.Debug("no-new-versions")
		return nil
	}

	logger.Info("versions-found", lager.Data{
		"versions": newVersions,
		"total":    len(newVersions),
	})

	version := newVersions[len(newVersions)-1]
	err = savedResourceType.SaveVersion(version)
	if err != nil {
		logger.Error("failed-to-save-resource-type-version", err, lager.Data{
			"version": version,
		})
		return err
	}

	return nil
}
