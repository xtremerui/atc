package pipes

import (
	"errors"
	"io"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
)

func (s *Server) WritePipe(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("write-pipe")
	pipeID := r.FormValue(":pipe_id")

	user, err := s.userFactory.GetUser(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", errors.New("failed-to-get-user"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbPipe, err := user.GetPipe(pipeID)
	if err != nil {
		logger.Error("failed-to-get-pipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.ID() != dbPipe.UserId {
		logger.Error("user-not-authorized-to-read-pipe",
			errors.New("user-not-authorized-to-read-pipe"),
			lager.Data{"UserId": user.ID(), "PipeID": dbPipe.ID})
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if dbPipe.URL == s.url {
		s.pipesL.RLock()
		pipe, found := s.pipes[pipeID]
		s.pipesL.RUnlock()

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = io.Copy(pipe.write, r.Body)
		_ = pipe.write.Close()

		s.pipesL.Lock()
		delete(s.pipes, pipeID)
		s.pipesL.Unlock()
	} else {
		logger.Debug("forwarding-pipe-write-request", lager.Data{"pipe-url": dbPipe.URL})
		response, err := s.forwardRequest(w, r, dbPipe.URL, atc.WritePipe, dbPipe.ID)
		if err != nil {
			logger.Error("failed-to-forward-request", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(response.StatusCode)
	}
}
