package pipelineserver

import (
	"net/http"

	"github.com/concourse/atc/api/accessor"
)

func (s *Server) UnpausePipeline(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("unpause-pipeline")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	pipeline, err := acc.TeamPipeline(accessor.Write, teamName, pipelineName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	err = pipeline.Unpause()
	if err != nil {
		logger.Error("failed-to-unpause-pipeline", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
