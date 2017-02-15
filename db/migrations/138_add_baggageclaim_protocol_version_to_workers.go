package migrations

import "github.com/concourse/atc/dbng/migration"

func AddBaggageclaimProtocolVersionToWorkers(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
	  ALTER TABLE workers
		ADD COLUMN baggageclaim_protocol_version int NOT NULL DEFAULT 0;
	`)
	if err != nil {
		return err
	}

	return nil
}
