package api

import (
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc/db"
)

type UserScopedHandlerFactory struct {
	logger      lager.Logger
	userFactory db.UserFactory
}

func NewUserScopedHandlerFactory(
	logger lager.Logger,
	userFactory db.UserFactory,
) *UserScopedHandlerFactory {
	return &UserScopedHandlerFactory{
		logger:      logger,
		userFactory: userFactory,
	}
}

func (f *UserScopedHandlerFactory) HandlerFor(userScopedHandler func(db.User) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := f.logger.Session("user-scoped-handler")

		user, err := f.userFactory.GetUser(r.Context())
		if err != nil {
			logger.Error("failed-to-create-user-from-context", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userScopedHandler(user).ServeHTTP(w, r)
	}
}
