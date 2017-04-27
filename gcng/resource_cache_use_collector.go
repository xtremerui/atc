package gcng

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
)

type resourceCacheUseCollector struct {
	logger       lager.Logger
	cacheFactory dbng.ResourceCacheFactory
}

func NewResourceCacheUseCollector(
	logger lager.Logger,
	cacheFactory dbng.ResourceCacheFactory,
) Collector {
	return &resourceCacheUseCollector{
		logger:       logger,
		cacheFactory: cacheFactory,
	}
}

func (rcuc *resourceCacheUseCollector) Run() error {
	logger := rcuc.logger.Session("run")

	logger.Debug("start")
	defer logger.Debug("done")

	err := rcuc.cacheFactory.CleanUsesForFinishedBuilds(logger.Session("for-finished-builds"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-builds", err)
		return err
	}

	err = rcuc.cacheFactory.CleanUsesForInactiveResourceTypes(logger.Session("for-inactive-resource-types"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-types", err)
		return err
	}

	err = rcuc.cacheFactory.CleanUsesForInactiveResources(logger.Session("for-inactive-resources"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-resources", err)
		return err
	}

	err = rcuc.cacheFactory.CleanUsesForPausedPipelineResources(logger.Session("for-paused-pipeline-resources"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-paused-pipeline-resources", err)
		return err
	}

	return nil
}
