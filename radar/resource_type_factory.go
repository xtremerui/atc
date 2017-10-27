package radar

import (
	"log"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
)

type resourceTypeFactory struct {
	dbPipeline db.Pipeline
	variables  creds.Variables
}

type ResourceTypeFactory interface {
	ResourceTypes(
		logger lager.Logger,
		resource string,
	) (atc.VersionedResourceTypes, error)
}

func NewResourceTypeFactory(
	dbPipeline db.Pipeline,
	variables creds.Variables,
) ResourceTypeFactory {
	return &resourceTypeFactory{
		dbPipeline: dbPipeline,
		variables:  variables,
	}
}

// ResourceTypes returns a list of resourceTypes that a given resource depends on
// sorted in reverse order i.e. furthest item (base type) on the dependency tree first
func (rtFactory *resourceTypeFactory) ResourceTypes(
	logger lager.Logger,
	resourceType string,
) (atc.VersionedResourceTypes, error) {
	logger.Session("resource-type-dependency-tree")
	var versionedResourceTypes atc.VersionedResourceTypes

	logger.Debug("Entering-Resource-Type-Factory", lager.Data{"resourceType": resourceType})

	pipelineResourceTypes, err := rtFactory.dbPipeline.ResourceTypes()
	if err != nil {
		logger.Error("failed-to-get-resource-types", err)
		return nil, err
	}

	credsPipelineResourceTypes := creds.NewVersionedResourceTypes(rtFactory.variables, pipelineResourceTypes.Deserialize())

	found := true
	lookupName := resourceType
	var customType creds.VersionedResourceType

	for found {
		customType, found = credsPipelineResourceTypes.Lookup(lookupName)
		if found {
			log.Println("FOUND", customType, found)

			credsPipelineResourceTypes = credsPipelineResourceTypes.Without(lookupName)

			lookupName = customType.Type
			versionedResourceTypes = append(atc.VersionedResourceTypes{customType.VersionedResourceType}, versionedResourceTypes...)
		}
	}

	return versionedResourceTypes, nil
}
