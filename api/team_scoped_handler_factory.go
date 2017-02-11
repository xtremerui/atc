package api

import (
	"errors"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
)

type TeamScopedHandlerFactory struct {
	logger        lager.Logger
	teamDBFactory db.TeamDBFactory
}

func NewTeamScopedHandlerFactory(
	logger lager.Logger,
	teamDBFactory db.TeamDBFactory,
) *TeamScopedHandlerFactory {
	return &TeamScopedHandlerFactory{
		logger:        logger,
		teamDBFactory: teamDBFactory,
	}
}

func (f *TeamScopedHandlerFactory) HandlerFor(teamScopedHandler func(db.TeamDB) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := f.logger.Session("team-scoped-handler")

		authTeam, authTeamFound := auth.GetTeam(r)
		if !authTeamFound {
			logger.Error("team-not-found-in-context", errors.New("team-not-found-in-context"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		teamDB, found, err := f.teamDBFactory.GetTeamDBByName(authTeam.Name())
		if err != nil {
			logger.Error("failed-to-get-team", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !found {
			logger.Debug("team-not-found", lager.Data{"team-name": authTeam.Name()})
			w.WriteHeader(http.StatusNotFound)
			return
		}

		teamScopedHandler(teamDB).ServeHTTP(w, r)
	}
}
