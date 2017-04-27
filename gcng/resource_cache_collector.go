package gcng

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
)

type resourceCacheCollector struct {
	logger       lager.Logger
	cacheFactory dbng.ResourceCacheFactory
}

func NewResourceCacheCollector(
	logger lager.Logger,
	cacheFactory dbng.ResourceCacheFactory,
) Collector {
	return &resourceCacheCollector{
		logger:       logger,
		cacheFactory: cacheFactory,
	}
}

func (rcuc *resourceCacheCollector) Run() error {
	logger := rcuc.logger.Session("run")
	logger.Debug("start")
	defer logger.Debug("done")

	return rcuc.cacheFactory.CleanUpInvalidCaches(logger)
}
