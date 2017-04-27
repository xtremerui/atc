package gcng

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
)

type resourceConfigUseCollector struct {
	logger        lager.Logger
	configFactory dbng.ResourceConfigFactory
}

func NewResourceConfigUseCollector(
	logger lager.Logger,
	configFactory dbng.ResourceConfigFactory,
) Collector {
	return &resourceConfigUseCollector{
		logger:        logger,
		configFactory: configFactory,
	}
}

func (rcuc *resourceConfigUseCollector) Run() error {
	logger := rcuc.logger.Session("run")

	logger.Debug("start")
	defer logger.Debug("done")

	err := rcuc.configFactory.CleanConfigUsesForFinishedBuilds(logger.Session("for-finished-builds"))
	if err != nil {
		logger.Error("unable-to-clean-up-config-uses", err)
		return err
	}

	err = rcuc.configFactory.CleanConfigUsesForInactiveResourceTypes(logger.Session("for-inactive-resource-types"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-types", err)
		return err
	}

	err = rcuc.configFactory.CleanConfigUsesForInactiveResources(logger.Session("for-inactive-resources"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-inactive-resources", err)
		return err
	}

	err = rcuc.configFactory.CleanConfigUsesForPausedPipelinesResources(logger.Session("for-paused-pipelines-resources"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-paused-resources", err)
		return err
	}

	err = rcuc.configFactory.CleanConfigUsesForOutdatedResourceConfigs(logger.Session("for-outdated-resource-configs"))
	if err != nil {
		logger.Error("unable-to-clean-up-for-outdated-resource-configs", err)
		return err
	}

	return nil
}
