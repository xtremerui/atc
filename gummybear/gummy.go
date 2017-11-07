package gummybear

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

type CheckRequest struct {
}

var resourceTypes []creds.VersionedResourceType
var variables creds.Variables

func CheckForVersions(
	logger lager.Logger,
	resourceFactory resource.ResourceFactory,
	resourceConfig db.Resource,
	workerSpec worker.WorkerSpec,
	owner db.ContainerOwner,
) ([]atc.Version, error) {

	metadata := resource.TrackerMetadata{
		ResourceName: resourceConfig.Name(),
		PipelineName: resourceConfig.PipelineName(),
	}

	containerSpec := worker.ContainerSpec{
		ImageSpec: worker.ImageSpec{
			ResourceType: resourceConfig.Type(),
		},
		Tags:   workerSpec.Tags,
		TeamID: workerSpec.TeamID,
		Env:    metadata.Env(),
	}

	source, err := creds.NewSource(variables, resourceConfig.Source()).Evaluate()
	if err != nil {
		logger.Error("failed-to-evaluate-resource-source", err)
		return []atc.Version{}, err
	}

	checkingResource, err := resourceFactory.NewResource(
		logger,
		nil,
		owner,
		db.ContainerMetadata{
			Type: db.ContainerTypeCheck,
		},
		containerSpec,
		resourceTypes,
		worker.NoopImageFetchingDelegate{},
	)
	if err != nil {
		return nil, err
	}

	return checkingResource.Check(source, nil)
}
