package auth

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/markbates/goth/gothic"
)

type LogOutHandler struct {
	logger lager.Logger
}

func NewLogOutHandler(logger lager.Logger) http.Handler {
	return &LogOutHandler{logger: logger}
}

func (handler *LogOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.logger.Session("logout")

	gothic.Logout(w, r)

	http.SetCookie(w, &http.Cookie{
		Name:   AuthCookieName,
		Path:   "/",
		MaxAge: -1,
	})
}
