package gcng

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
)

type resourceConfigCollector struct {
	logger        lager.Logger
	configFactory dbng.ResourceConfigFactory
}

func NewResourceConfigCollector(
	logger lager.Logger,
	configFactory dbng.ResourceConfigFactory,
) Collector {
	return &resourceConfigCollector{
		logger:        logger,
		configFactory: configFactory,
	}
}

func (rcuc *resourceConfigCollector) Run() error {
	logger := rcuc.logger.Session("run")
	logger.Debug("start")
	defer logger.Debug("done")

	return rcuc.configFactory.CleanUselessConfigs(logger)
}
