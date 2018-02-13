package auth

import (
	"net/http"

	"code.cloudfoundry.org/lager"
)

func CSRFValidationHandler(
	handler http.Handler,
	rejector Rejector,
	tokenValidator TokenValidator,
) http.Handler {
	return csrfValidationHandler{
		handler:        handler,
		rejector:       rejector,
		tokenValidator: tokenValidator,
	}
}

type csrfValidationHandler struct {
	handler        http.Handler
	rejector       Rejector
	tokenValidator TokenValidator
}

func (h csrfValidationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger, ok := r.Context().Value("logger").(lager.Logger)
	if !ok {
		panic("logger is not set in request context for csrf validation handler")
	}

	logger = logger.Session("csrf-validation")

	// We don't validate CSRF token for GET requests
	// since they are not changing the state
	if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
		isCSRFRequired, ok := r.Context().Value(CSRFRequiredKey).(bool)
		if ok && isCSRFRequired {
			if r.Header.Get(CSRFHeaderName) == "" {
				logger.Debug("csrf-header-is-not-set")
				h.rejector.Unauthorized(w, r)
				return
			}

			authCSRFToken, authCSRFTokenProvided := h.tokenValidator.GetCSRFToken(r)
			if !authCSRFTokenProvided {
				logger.Debug("csrf-is-not-provided-in-auth-token")
				h.rejector.Unauthorized(w, r)
				return
			}

			if authCSRFToken != r.Header.Get(CSRFHeaderName) {
				logger.Debug("csrf-token-does-not-match-auth-token", lager.Data{
					"auth-csrf-token":    authCSRFToken,
					"request-csrf-token": r.Header.Get(CSRFHeaderName),
				})
				h.rejector.Unauthorized(w, r)
				return
			}
		}
	}

	h.handler.ServeHTTP(w, r)
}
