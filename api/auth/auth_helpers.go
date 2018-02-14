package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func IsAdmin(r *http.Request) bool {
	isAdmin, present := r.Context().Value(isAdminKey).(bool)
	return present && isAdmin
}

func IsSystem(r *http.Request) bool {
	isSystem, present := r.Context().Value(isSystemKey).(bool)
	return present && isSystem
}

func IsAuthenticated(r *http.Request) bool {
	isAuthenticated, _ := r.Context().Value(isAuthenticatedKey).(bool)
	return isAuthenticated
}

func IsAuthorized(r *http.Request) bool {
	authorizer, authFound := GetAuthorizer(r)

	if authFound && authorizer.IsAuthorized(r.URL.Query().Get(":team_name")) {
		return true
	}

	return false
}

func getJWT(r *http.Request, publicKey *rsa.PublicKey) (token *jwt.Token, err error) {
	fun := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	}

	if ah := r.Header.Get("Authorization"); ah != "" {
		// Should be a bearer token
		if len(ah) > 6 && strings.ToUpper(ah[0:6]) == "BEARER" {
			return jwt.Parse(ah[7:], fun)
		}
	}

	return nil, errors.New("unable to parse authorization header")
}

func GetAuthorizer(r *http.Request) (Authorizer, bool) {
	userId, userIdPresent := r.Context().Value(userIdKey).(string)
	teams, teamsPresent := r.Context().Value(teamsKey).([]string)
	isAdmin, adminPresent := r.Context().Value(isAdminKey).(bool)

	if !(userIdPresent && teamsPresent && adminPresent) {
		return nil, false
	}

	return &authorizer{userId, teams, isAdmin}, true
}

type Authorizer interface {
	UserId() string
	Teams() []string
	IsAdmin() bool
	IsAuthorized(teamName string) bool
}

type authorizer struct {
	userId  string
	teams   []string
	isAdmin bool
}

func (a *authorizer) UserId() string {
	return a.userId
}

func (a *authorizer) Teams() []string {
	return a.teams
}

func (a *authorizer) IsAdmin() bool {
	return a.isAdmin
}

func (a *authorizer) IsAuthorized(teamName string) bool {
	for _, team := range a.teams {
		if team == teamName {
			return true
		}
	}
	return false
}
