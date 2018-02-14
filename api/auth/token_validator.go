package auth

import (
	"crypto/rsa"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

//go:generate counterfeiter . TokenValidator

type TokenValidator interface {
	IsAuthenticated(r *http.Request) bool
	GetUserId(r *http.Request) (string, bool)
	GetTeams(r *http.Request) ([]string, bool, bool)
	GetSystem(r *http.Request) (bool, bool)
	GetCSRFToken(r *http.Request) (string, bool)
}

type JWTValidator struct {
	PublicKey *rsa.PublicKey
}

func (v JWTValidator) IsAuthenticated(r *http.Request) bool {
	token, err := getJWT(r, v.PublicKey)
	if err != nil {
		return false
	}

	return token.Valid
}

func (v JWTValidator) GetUserId(r *http.Request) (string, bool) {
	token, err := getJWT(r, v.PublicKey)
	if err != nil {
		return "", false
	}

	claims := token.Claims.(jwt.MapClaims)
	userIdInterface, userIdOK := claims[userIdClaimKey]

	if !(userIdOK) {
		return "", false
	}

	userId := userIdInterface.(string)

	return userId, true
}

func (v JWTValidator) GetTeams(r *http.Request) ([]string, bool, bool) {
	token, err := getJWT(r, v.PublicKey)
	if err != nil {
		return []string{}, false, false
	}

	claims := token.Claims.(jwt.MapClaims)
	teamsInterface, teamNameOK := claims[teamsClaimKey]
	isAdminInterface, isAdminOK := claims[isAdminClaimKey]

	if !(teamNameOK && isAdminOK) {
		return []string{}, false, false
	}

	teams := teamsInterface.([]string)
	isAdmin := isAdminInterface.(bool)

	return teams, isAdmin, true
}

func (v JWTValidator) GetSystem(r *http.Request) (bool, bool) {
	token, err := getJWT(r, v.PublicKey)
	if err != nil {
		return false, false
	}

	claims := token.Claims.(jwt.MapClaims)
	isSystemInterface, isSystemOK := claims[isSystemKey]
	if !isSystemOK {
		return false, false
	}

	return isSystemInterface.(bool), true
}

func (v JWTValidator) GetCSRFToken(r *http.Request) (string, bool) {
	token, err := getJWT(r, v.PublicKey)
	if err != nil {
		return "", false
	}

	claims := token.Claims.(jwt.MapClaims)
	csrfToken, ok := claims[csrfTokenClaimKey]
	if !ok {
		return "", false
	}

	return csrfToken.(string), true
}
