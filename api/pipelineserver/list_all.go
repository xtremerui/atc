package pipelineserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
)

// show all public pipelines and team private pipelines if authorized
func (s *Server) ListAllPipelines(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-all-pipelines")
	authTeam, authTeamFound := auth.GetTeam(r)

	var pipelines []db.SavedPipeline
	var err error
	var teamDB db.TeamDB
	if authTeamFound {
		var found bool
		teamDB, found, err = s.teamDBFactory.GetTeamDBByName(authTeam.Name())
		if err != nil {
			logger.Error("failed-to-get-all-team", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !found {
			logger.Debug("team-not-found", lager.Data{"team-name": authTeam.Name()})
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pipelines, err = teamDB.GetPrivateAndAllPublicPipelines()
	} else {
		pipelines, err = s.pipelinesDB.GetAllPublicPipelines()
	}

	if err != nil {
		logger.Error("failed-to-get-all-active-pipelines", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(present.Pipelines(pipelines))
}
