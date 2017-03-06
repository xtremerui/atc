package migrations

import (
	"fmt"

	"github.com/concourse/atc/dbng/migration"
)

func AddWorkerResourceCacheToVolumes(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
      ALTER TABLE volumes
      ADD COLUMN worker_resource_cache_id INTEGER
  		REFERENCES worker_resource_caches (id) ON DELETE SET NULL
		`)
	if err != nil {
		return err
	}

	rows, err := tx.Query(`SELECT id, resource_cache_id, worker_base_resource_type_id FROM volumes WHERE resource_cache_id IS NOT NULL`)
	if err != nil {
		return err
	}

	defer rows.Close()

	volumeWorkerResourceCaches := []volumeWorkerResourceCache{}

	for rows.Next() {
		var id int
		var resourceCacheID int
		var workerBaseResourceTypeID int
		err = rows.Scan(&id, &resourceCacheID, &workerBaseResourceTypeID)
		if err != nil {
			return fmt.Errorf("failed to scan volume id, resource_cache_id and worker_name: %s", err)
		}

		volumeWorkerResourceCaches = append(volumeWorkerResourceCaches, volumeWorkerResourceCache{
			ID:                       id,
			ResourceCacheID:          resourceCacheID,
			WorkerBaseResourceTypeID: workerBaseResourceTypeID,
		})
	}

	for _, vwrc := range volumeWorkerResourceCaches {
		var workerResourceCacheID int
		err = tx.QueryRow(`
				INSERT INTO worker_resource_caches (worker_base_resource_type_id, resource_cache_id)
		    VALUES ($1, $2)
		    RETURNING id
			`, vwrc.WorkerBaseResourceTypeID, vwrc.ResourceCacheID).
			Scan(&workerResourceCacheID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
        UPDATE volumes SET worker_resource_cache_id=$1 WHERE id=$2
      `, workerResourceCacheID, vwrc.ID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
      ALTER TABLE containers
      DROP COLUMN resource_cache_id
      DROP COLUMN worker_base_resource_type_id
    `)
	if err != nil {
		return err
	}

	return nil
}

type volumeWorkerResourceCache struct {
	ID                       int
	ResourceCacheID          int
	WorkerBaseResourceTypeID int
}
