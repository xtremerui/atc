package db_test

import (
	"context"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Factory", func() {

	Describe("Workers", func() {
		var (
			user      db.User
			team1     db.Team
			team2     db.Team
			atcWorker atc.Worker
			err       error
		)

		BeforeEach(func() {
			postgresRunner.Truncate()
			team1, err = teamFactory.CreateTeam(atc.Team{Name: "team"})
			Expect(err).ToNot(HaveOccurred())

			team2, err = teamFactory.CreateTeam(atc.Team{Name: "some-other-team"})
			Expect(err).ToNot(HaveOccurred())

			atcWorker = atc.Worker{
				GardenAddr:       "some-garden-addr",
				BaggageclaimURL:  "some-bc-url",
				HTTPProxyURL:     "some-http-proxy-url",
				HTTPSProxyURL:    "some-https-proxy-url",
				NoProxy:          "some-no-proxy",
				ActiveContainers: 140,
				ResourceTypes: []atc.WorkerResourceType{
					{
						Type:    "some-resource-type",
						Image:   "some-image",
						Version: "some-version",
					},
					{
						Type:    "other-resource-type",
						Image:   "other-image",
						Version: "other-version",
					},
				},
				Platform:  "some-platform",
				Tags:      atc.Tags{"some", "tags"},
				Name:      "some-name",
				StartTime: 55,
			}

			c := context.Background()
			c = context.WithValue(c, "teams", []string{team1.Name()})
			user, err = userFactory.GetUser(c)
		})

		Context("when there are global workers and workers for the team", func() {
			BeforeEach(func() {
				_, err = team1.SaveWorker(atcWorker, 0)
				Expect(err).ToNot(HaveOccurred())

				atcWorker.Name = "some-new-worker"
				atcWorker.GardenAddr = "some-other-garden-addr"
				atcWorker.BaggageclaimURL = "some-other-bc-url"
				_, err = workerFactory.SaveWorker(atcWorker, 0)
				Expect(err).ToNot(HaveOccurred())
			})

			It("finds them without error", func() {
				workers, err := user.Workers()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(workers)).To(Equal(2))

				Expect(workers[0].Name()).To(Equal("some-name"))
				Expect(*workers[0].GardenAddr()).To(Equal("some-garden-addr"))
				Expect(*workers[0].BaggageclaimURL()).To(Equal("some-bc-url"))

				Expect(workers[1].Name()).To(Equal("some-new-worker"))
				Expect(*workers[1].GardenAddr()).To(Equal("some-other-garden-addr"))
				Expect(*workers[1].BaggageclaimURL()).To(Equal("some-other-bc-url"))
			})
		})

		Context("when there are workers for another team", func() {
			BeforeEach(func() {
				atcWorker.Name = "some-other-team-worker"
				atcWorker.GardenAddr = "some-other-garden-addr"
				atcWorker.BaggageclaimURL = "some-other-bc-url"
				_, err = team2.SaveWorker(atcWorker, 5*time.Minute)
				Expect(err).ToNot(HaveOccurred())
			})

			It("does not find the other team workers", func() {
				workers, err := user.Workers()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(workers)).To(Equal(0))
			})
		})
	})
})
