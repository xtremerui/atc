package db_test

import (
	"time"

	"github.com/concourse/atc/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("ContainerOwner", func() {

	Describe("ResourceConfigCheckSessionContainerOwner", func() {
		var (
			owner         db.ContainerOwner
			ownerExpiries db.ContainerOwnerExpiries
			found         bool
		)

		BeforeEach(func() {



			resourceConfigFactory.FindOrCreateResourceConfig(logger, 
		
		
	  resourceConfigUser,	
		resourceType string,
		source atc.Source,
		resourceTypes creds.VersionedResourceTypes,
		)

			resourceConfig := &db.UsedResourceConfig{
				ID: 1,
				CreatedByBaseResourceType: &db.UsedBaseResourceType{
					ID:   1,
					Name: "fake-resource-type",
				},
			}

			ownerExpiries = db.ContainerOwnerExpiries{
				ExpiryGraceTime: 1 * time.Minute,
				MinExpiry:       1 * time.Minute,
				MaxExpiry:       1 * time.Minute,
			}

			owner = db.NewResourceConfigCheckSessionContainerOwner(
				resourceConfig,
				ownerExpiries,
			)
		})

		Describe("Find", func() {
			JustBeforeEach(func() {
				var err error
				_, found, err = owner.Find(dbConn)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("When a resource config check session exists", func() {
				Context("and the ExpiryGraceTime is 1 minute", func() {
					Context("and it will not expire in the next minute", func() {
						BeforeEach(func() {
							_, err := psql.Insert("resource_config_check_sessions").
								SetMap(map[string]interface{}{
									"resource_config_id":           1,
									"worker_base_resource_type_id": 2,
									"expires_at":                   time.Now().Add(10 * time.Minute),
								}).
								RunWith(dbConn).
								Exec()

							Expect(err).ToNot(HaveOccurred())
						})

						It("finds the resource config check session", func() {
							Expect(found).To(BeTrue())
						})
					})

					Context("and it expires in less than a minute", func() {
						BeforeEach(func() {
							//make a rccs that doesn't expire in 1 minute (time.Now -1 minute )
							_, err := psql.Insert("resource_config_check_sessions").
								SetMap(map[string]interface{}{
									"resource_config_id":           1,
									"worker_base_resource_type_id": 2,
									"expires_at":                   time.Now().Add(-1 * time.Minute),
								}).
								RunWith(dbConn).
								Exec()
							Expect(err).ToNot(HaveOccurred())
						})

						It("doesn't find a resource config check session", func() {
							Expect(found).To(BeFalse())
						})
					})
				})
			})

			Context("When a resource config check session doesn't exist", func() {
				It("doesn't find a resource config check session", func() {
					Expect(found).To(BeFalse())
				})
			})
		})

		Describe("Create", func() {

		})

	})
})
