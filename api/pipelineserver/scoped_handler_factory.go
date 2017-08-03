package pipelineserver

import (
	"fmt"
	"net/http"

	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
	"github.com/google/jsonapi"
)

type ScopedHandlerFactory struct {
	teamDBFactory db.TeamFactory
}

func NewScopedHandlerFactory(
	teamDBFactory db.TeamFactory,
) *ScopedHandlerFactory {
	return &ScopedHandlerFactory{
		teamDBFactory: teamDBFactory,
	}
}

func (pdbh *ScopedHandlerFactory) HandlerFor(pipelineScopedHandler func(db.Pipeline) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.FormValue(":team_name")
		pipelineName := r.FormValue(":pipeline_name")

		pipeline, ok := r.Context().Value(auth.PipelineContextKey).(db.Pipeline)
		if !ok {
			dbTeam, found, err := pdbh.teamDBFactory.FindTeam(teamName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !found {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
					Title:  "Team Not Found Error",
					Detail: fmt.Sprintf("Team with name '%s' not found.", teamName),
					Status: "404",
				}})
				return
			}

			pipeline, found, err = dbTeam.Pipeline(pipelineName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !found {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
					Title:  "Pipeline Not Found Error",
					Detail: fmt.Sprintf("Pipeline with name '%s' not found.", pipelineName),
					Status: "404",
				}})
				return
			}
		}

		pipelineScopedHandler(pipeline).ServeHTTP(w, r)
	}
}
