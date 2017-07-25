package integration_test

import (
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/postgresrunner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"

	"testing"

	rfakes "github.com/concourse/atc/resource/resourcefakes"
)

var (
	postgresRunner postgresrunner.Runner
	dbProcess      ifrit.Process

	dbConn                              db.Conn
	lockFactory                         lock.LockFactory
	teamFactory                         db.TeamFactory
	workerFactory                       db.WorkerFactory
	resourceConfigCheckSessionLifecycle db.ResourceConfigCheckSessionLifecycle
	resourceConfigFactory               db.ResourceConfigFactory
	fakeResourceFactory                 *rfakes.FakeResourceFactory
	variables                           creds.Variables

	defaultWorkerResourceType atc.WorkerResourceType
	defaultTeam               db.Team
	defaultWorkerPayload      atc.Worker
	defaultWorker             db.Worker
	defaultResourceType       db.ResourceType
	defaultResource           db.Resource
	defaultPipeline           db.Pipeline
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	postgresRunner = postgresrunner.Runner{
		Port: 5433 + GinkgoParallelNode(),
	}

	dbProcess = ifrit.Invoke(postgresRunner)

	postgresRunner.CreateTestDB()
	dbConn = postgresRunner.OpenConn()
	lockFactory = lock.NewLockFactory(postgresRunner.OpenSingleton())
	teamFactory = db.NewTeamFactory(dbConn, lockFactory)
	variables = template.StaticVariables{
		"source-params": "some-secret-sauce",
	}

	resourceConfigCheckSessionLifecycle = db.NewResourceConfigCheckSessionLifecycle(dbConn)
	resourceConfigFactory = db.NewResourceConfigFactory(dbConn, lockFactory)
	fakeResourceFactory = new(rfakes.FakeResourceFactory)

	workerFactory = db.NewWorkerFactory(dbConn)

	var err error
	defaultTeam, err = teamFactory.CreateTeam(atc.Team{Name: "default-team"})
	Expect(err).NotTo(HaveOccurred())

	defaultWorkerResourceType = atc.WorkerResourceType{
		Type:    "some-base-resource-type",
		Image:   "/path/to/image",
		Version: "some-brt-version",
	}

	defaultWorkerPayload = atc.Worker{
		ResourceTypes:   []atc.WorkerResourceType{defaultWorkerResourceType},
		Name:            "default-worker",
		GardenAddr:      "1.2.3.4:7777",
		BaggageclaimURL: "5.6.7.8:7878",
	}

	defaultWorker, err = workerFactory.SaveWorker(defaultWorkerPayload, 0)
	Expect(err).NotTo(HaveOccurred())

	defaultPipeline, _, err = defaultTeam.SavePipeline("default-pipeline", atc.Config{
		Jobs: atc.JobConfigs{
			{
				Name: "some-job",
			},
		},
		Resources: atc.ResourceConfigs{
			{
				Name: "some-resource",
				Type: "some-base-resource-type",
				Source: atc.Source{
					"some": "source",
				},
			},
		},
		ResourceTypes: atc.ResourceTypes{
			{
				Name: "some-type",
				Type: "some-base-resource-type",
				Source: atc.Source{
					"some-type": "source",
				},
			},
		},
	}, db.ConfigVersion(0), db.PipelineUnpaused)
	Expect(err).NotTo(HaveOccurred())

	var found bool
	defaultResourceType, found, err = defaultPipeline.ResourceType("some-type")
	Expect(found).To(BeTrue())
	Expect(err).NotTo(HaveOccurred())

	err = defaultResourceType.SaveVersion(atc.Version{"some-type": "version"})
	Expect(err).NotTo(HaveOccurred())

	found, err = defaultResourceType.Reload()
	Expect(err).NotTo(HaveOccurred())
	Expect(found).To(BeTrue())

	defaultResource, found, err = defaultPipeline.Resource("some-resource")
	Expect(err).NotTo(HaveOccurred())
	Expect(found).To(BeTrue())
})

var _ = AfterEach(func() {
	postgresRunner.Truncate()

})
