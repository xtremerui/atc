package migrations

import "github.com/concourse/atc/db/migration"

func CreateJobVariants(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
		CREATE TABLE job_variants (
			id serial PRIMARY KEY,
			job_id int REFERENCES jobs (id) ON DELETE CASCADE,
			active bool DEFAULT true NOT NULL,
			inputs_determined bool DEFAULT false NOT NULL,
			resource_spaces jsonb NOT NULL,
			UNIQUE (job_id, resource_spaces)
		)
	`)
	return err
}
