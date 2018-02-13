package auth

import (
	"context"
	"net/http"
)

func WrapHandler(
	handler http.Handler,
	validator TokenValidator,
) http.Handler {
	return authHandler{
		handler:   handler,
		validator: validator,
	}
}

type authHandler struct {
	handler   http.Handler
	validator TokenValidator
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	isAuthenticated := h.validator.IsAuthenticated(r)
	ctx := context.WithValue(r.Context(), authenticated, isAuthenticated)

	teamName, isAdmin, found := h.validator.GetTeam(r)
	if found {
		ctx = context.WithValue(ctx, teamNameKey, teamName)
		ctx = context.WithValue(ctx, isAdminKey, isAdmin)
	}

	isSystem, found := h.validator.GetSystem(r)
	if found {
		ctx = context.WithValue(ctx, isSystemKey, isSystem)
	}

	h.handler.ServeHTTP(w, r.WithContext(ctx))
}
