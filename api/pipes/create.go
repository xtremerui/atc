package pipes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/tedsuo/rata"

	"github.com/concourse/atc"
)

func (s *Server) CreatePipe(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("create-pipe")

	user, err := s.userFactory.GetUser(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", errors.New("failed-to-get-user"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	teamName := r.FormValue(":team_name")
	team, authorized, err := user.GetTeam(teamName)
	if err != nil {
		logger.Error("failed-to-get-team", errors.New("failed-to-get-team"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	guid, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = team.CreatePipe(guid.String(), s.url)
	if err != nil {
		logger.Error("failed-to-create-pipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pr, pw := io.Pipe()

	pipeID := guid.String()

	reqGen := rata.NewRequestGenerator(s.externalURL, atc.Routes)

	readReq, err := reqGen.CreateRequest(atc.ReadPipe, rata.Params{
		"pipe_id": pipeID,
	}, nil)
	if err != nil {
		logger.Error("failed-to-create-pipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeReq, err := reqGen.CreateRequest(atc.WritePipe, rata.Params{
		"pipe_id": pipeID,
	}, nil)
	if err != nil {
		logger.Error("failed-to-create-pipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pipeResource := atc.Pipe{
		ID:       pipeID,
		ReadURL:  readReq.URL.String(),
		WriteURL: writeReq.URL.String(),
	}

	pipe := pipe{
		resource: pipeResource,

		read:  pr,
		write: pw,
	}

	s.pipesL.Lock()
	s.pipes[pipeResource.ID] = pipe
	s.pipesL.Unlock()

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(pipeResource)
	if err != nil {
		logger.Error("failed-to-encode-pipe", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
