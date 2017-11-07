package migration_test

import (
	"database/sql"
	"fmt"

	"github.com/concourse/atc/db/migration"
	"github.com/concourse/atc/db/migrations"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateJobVariants", func() {
	var dbConn *sql.DB
	var migrator migration.Migrator

	// explicit type here is important for reflect.ValueOf
	migrator = migrations.CreateJobVariants

	BeforeEach(func() {
		var err error
		dbConn, err = openDBConnPreMigration(migrator)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := dbConn.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	FContext("before job variants are introduced", func() {
		var teamID, pipelineID, jobID, resourceID, resourceSpaceID int

		BeforeEach(func() {
			err := dbConn.QueryRow(`
				INSERT INTO teams(name) VALUES('some-team') RETURNING id
			`).Scan(&teamID)
			Expect(err).NotTo(HaveOccurred())

			err = dbConn.QueryRow(`
				INSERT INTO pipelines(name, team_id) VALUES('some-pipeline', $1) RETURNING id
			`, teamID).Scan(&pipelineID)
			Expect(err).NotTo(HaveOccurred())

			err = dbConn.QueryRow(`
				INSERT INTO resources(name, config, pipeline_id) VALUES('some-resource', $1, $2) RETURNING id
			`, "{}", pipelineID).Scan(&resourceID)
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("resourceID %d", resourceID)
			err = dbConn.QueryRow(`
				INSERT INTO resource_spaces(name, resource_id) VALUES('some-space', $1) RETURNING id
			`, resourceID).Scan(&resourceSpaceID)
			fmt.Printf("resourceSpaceID %d", resourceSpaceID)
			Expect(err).NotTo(HaveOccurred())

			err = dbConn.QueryRow(`
				INSERT INTO jobs(name, pipeline_id, config, active, inputs_determined) VALUES('some-job', $1, '{}', true, true) RETURNING id
			`, pipelineID).Scan(&jobID)
			Expect(err).NotTo(HaveOccurred())

			err = dbConn.Close()
			Expect(err).NotTo(HaveOccurred())

			dbConn, err = openDBConnPostMigration(migrator)
			Expect(err).NotTo(HaveOccurred())
		})

		It("migrates the jobs to job variants", func() {
			var jobVariantID int
			var active, inputsDetermined bool
			var resourceSpaces string
			err := dbConn.QueryRow(`
				SELECT id, active, inputs_determined, resource_spaces FROM job_variants WHERE job_id=$1
			`, jobID).Scan(&jobVariantID, &active, &inputsDetermined, &resourceSpaces)
			Expect(err).NotTo(HaveOccurred())

			Expect(jobVariantID).To(Equal(jobID))
			Expect(active).To(BeTrue())
			Expect(inputsDetermined).To(BeTrue())
			Expect(resourceSpaces).To(Equal("{\"some-resource\": \"some-space\"}"))
		})
	})
})
