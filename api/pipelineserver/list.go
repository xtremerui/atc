package pipelineserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
)

func (s *Server) ListPipelines(w http.ResponseWriter, r *http.Request) {
	var pipelines []db.SavedPipeline
	var err error

	logger := s.logger.Session("list-pipelines")
	requestTeamName := r.FormValue(":team_name")
	teamDB, found, err := s.teamDBFactory.GetTeamDBByName(requestTeamName)
	if err != nil {
		logger.Error("failed-to-get-team", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		logger.Debug("team-not-found", lager.Data{"team-name": requestTeamName})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	authTeam, authTeamFound := auth.GetTeam(r)
	if authTeamFound && authTeam.IsAuthorized(requestTeamName) {
		pipelines, err = teamDB.GetPipelines()
	} else {
		pipelines, err = teamDB.GetPublicPipelines()
	}

	if err != nil {
		logger.Error("failed-to-get-all-active-pipelines", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(present.Pipelines(pipelines))
}
