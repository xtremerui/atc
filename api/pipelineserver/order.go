package pipelineserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
)

func (s *Server) OrderPipelines(w http.ResponseWriter, r *http.Request) {
	pipelineNames := []string{}

	if err := json.NewDecoder(r.Body).Decode(&pipelineNames); err != nil {
		s.logger.Error("invalid-json", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	teamName := r.FormValue(":team_name")
	teamDB, found, err := s.teamDBFactory.GetTeamDBByName(teamName)
	if err != nil {
		s.logger.Error("failed-to-get-team", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		s.logger.Debug("team-not-found", lager.Data{"team-name": teamName})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = teamDB.OrderPipelines(pipelineNames)
	if err != nil {
		s.logger.Error("failed-to-order-pipelines", err, lager.Data{
			"pipeline-names": pipelineNames,
		})

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
