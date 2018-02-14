package workerserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/db"
)

func (s *Server) ListWorkers(user db.User) http.Handler {
	logger := s.logger.Session("list-workers")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		savedWorkers, err := user.Workers()
		if err != nil {
			logger.Error("failed-to-get-workers", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		workers := make([]atc.Worker, len(savedWorkers))
		for i, savedWorker := range savedWorkers {
			workers[i] = present.Worker(savedWorker)
		}

		err = json.NewEncoder(w).Encode(workers)
		if err != nil {
			logger.Error("failed-to-encode-workers", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
