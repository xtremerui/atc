package integration_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/gc"
	"github.com/concourse/atc/gc/gcfakes"
	"github.com/concourse/atc/radar"
	rfakes "github.com/concourse/atc/resource/resourcefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource Checking Container Garbage Collection", func() {

	var (
		gcLogger    *lagertest.TestLogger
		radarLogger *lagertest.TestLogger

		fakeClock *fakeclock.FakeClock
		interval  time.Duration

		containerFactory db.ContainerFactory

		fakeResource *rfakes.FakeResource

		fakeJobRunner      *gcfakes.FakeWorkerJobRunner
		containerCollector gc.Collector

		scanner radar.Scanner
	)

	BeforeEach(func() {
		gcLogger = lagertest.NewTestLogger("gc")
		radarLogger = lagertest.NewTestLogger("radar")
		containerFactory = db.NewContainerFactory(dbConn)

		fakeResource = new(rfakes.FakeResource)

		containerCollector = gc.NewContainerCollector(
			gcLogger,
			containerFactory,
			fakeJobRunner,
		)

		scanner = radar.NewResourceScanner(
			fakeClock,
			fakeResourceFactory,
			resourceConfigFactory,
			interval,
			defaultPipeline,
			"https://www.example.com",
			variables,
		)
	})

	Context("when radar is currently checking for new versions of a resource", func() {
		var (
			scannerErr   error
			collectorErr error
		)

		BeforeEach(func() {
			fakeResourceFactory.NewResourceReturns(fakeResource, nil)
			fakeResource.CheckReturns([]atc.Version{}, nil)
		})

		Context("and the Garbage Collector has run", func() {

			BeforeEach(func() {
				//hey radar make some checks and stuff
				go func() {
					_, scannerErr = scanner.Run(radarLogger, "some-resource")
				}()

				go func() {
					collectorErr = containerCollector.Run()
					fakeResourceFactory.NewResourceReturns(fakeResource, errors.New("/proc/123/ns/net blah"))
				}()
			})

			It("does not remove the resource checking container for the resource", func() {
				Consistently(func() error {
					return scannerErr
				}).ShouldNot(HaveOccurred())

				Consistently(func() error {
					return collectorErr
				}).ShouldNot(HaveOccurred())
			})
		})
	})

})
