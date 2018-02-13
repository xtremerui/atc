package auth

import (
	"crypto/rsa"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

//go:generate counterfeiter . TokenValidator

type TokenValidator interface {
	IsAuthenticated(r *http.Request) bool
	GetTeam(r *http.Request) (string, bool, bool)
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

func (v JWTValidator) GetTeam(r *http.Request) (string, bool, bool) {
	token, err := getJWT(r, v.PublicKey)
	if err != nil {
		return "", false, false
	}

	claims := token.Claims.(jwt.MapClaims)
	teamNameInterface, teamNameOK := claims[teamNameClaimKey]
	isAdminInterface, isAdminOK := claims[isAdminClaimKey]

	if !(teamNameOK && isAdminOK) {
		return "", false, false
	}

	teamName := teamNameInterface.(string)
	isAdmin := isAdminInterface.(bool)

	return teamName, isAdmin, true
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
