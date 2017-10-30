package radar

import (
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	. "github.com/concourse/atc/radar"

	rfakes "github.com/concourse/atc/resource/resourcefakes"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("ResourceTypeFactory", func() {
	var (
		fakeResourceFactory                   *rfakes.FakeResourceFactory
		fakeResourceConfigCheckSessionFactory *dbfakes.FakeResourceConfigCheckSessionFactory
		fakeResourceConfigCheckSession        *dbfakes.FakeResourceConfigCheckSession
		fakePipeline                          *dbfakes.FakePipeline
		variables                             creds.Variables

		fakeResourceType *dbfakes.FakeResourceType

		teamID = 123

		factory ResourceTypeFactory
	)

	BeforeEach(func() {
		variables = template.StaticVariables{
			"source-params": "some-secret-sauce",
		}

		fakeResourceFactory = new(rfakes.FakeResourceFactory)
		fakeResourceConfigCheckSessionFactory = new(dbfakes.FakeResourceConfigCheckSessionFactory)
		fakeResourceConfigCheckSession = new(dbfakes.FakeResourceConfigCheckSession)
		fakeResourceType = new(dbfakes.FakeResourceType)
		fakePipeline = new(dbfakes.FakePipeline)

		fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionReturns(fakeResourceConfigCheckSession, nil)

		fakeResourceConfigCheckSession.ResourceConfigReturns(&db.UsedResourceConfig{ID: 123})

		fakePipeline.IDReturns(42)
		fakePipeline.NameReturns("some-pipeline")
		fakePipeline.TeamIDReturns(teamID)
		fakePipeline.ReloadReturns(true, nil)
		fakePipeline.ResourceTypesReturns([]db.ResourceType{fakeResourceType}, nil)
		fakePipeline.ResourceTypeReturns(fakeResourceType, true, nil)

		factory = NewResourceTypeFactory(
			fakeResourceFactory,
			fakeResourceConfigCheckSessionFactory,
			fakePipeline,
			variables,
		)
	})

	Describe("DependentResourceTypes", func() {
		var types creds.VersionedResourceTypes

		BeforeEach(func() {
			givenType = "some-type"
			types = nil
		})

		JustBeforeEach(func() {
			var err error
			types, err = factory.DependentResourceTypes(givenType)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when given a base resource type", func() {
			BeforeEach(func() {
				fakePipeline.ResourceTypesReturns(db.ResourceTypes{}, nil)
			})
		})

		Context("when given a pipeline resource type", func() {
			var fakeSomeType *dbfakes.FakeResourceType

			BeforeEach(func() {
				fakeSomeType  =new(dbfakes.FakeResourceType)
				fakeSomeType.IDReturns(39)
				fakeSomeType.NameReturns("some-type")
				fakeSomeType.TypeReturns("docker-image")
				fakeSomeType.SourceReturns(atc.Source{"custom": "((source-params))"})
				fakeSomeType.SetResourceConfigReturns(nil)

				fakePipeline.ResourceTypesReturns(db.ResourceTypes{fakeResourceType}, nil)
			})

			Context("when the type has no version", func() {
				BeforeEach(func() {
					fakeSomeType.VersionReturns(nil)
				})
			})

			Context("when the type's version is present", func() {
				BeforeEach(func() {
					fakeSomeType.VersionReturns(atc.Version{"some": "version"})
				})
			})

			Context("when the type depends on another type", func() {
				var fakeSomeOtherType *dbfakes.FakeResourceType

				BeforeEach(func() {
					fakeSomeType.TypeReturns("some-other-type")

				fakeSomeOtherType  =new(dbfakes.FakeResourceType)
				fakeSomeOtherType.IDReturns(40)
				fakeSomeOtherType.NameReturns("some-other-type")
				fakeSomeOtherType.TypeReturns("docker-image")
				fakeSomeOtherType.SourceReturns(atc.Source{"custom": "((source-params))"})
				fakeSomeOtherType.SetResourceConfigReturns(nil)
				})

				Context("when the dependent type has no version", func() {
					BeforeEach(func() {
						fakeSomeOtherType.V
				})

				Context("when the dependent type's version is present", func() {
				})
			})
		})
	})
})
