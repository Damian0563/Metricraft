package main

import (
	"net/http"
	"strings"
)

func (t Token) GetUser() (User, error) {
	return getUserByToken(t.token)
}

func (t Token) GetAppName() (string, error) {
	split := strings.Split(t.token, ":")
	return getAppNameByToken(split[0])
}

func (t Token) sign(rotate bool) (string, error) {
	split := strings.Split(t.token, ":")
	return signToken(split[0], rotate)
}

func (t Token) verify() (bool, error) {
	return checkToken(t.token)
}

func (t Token) updateToken() error {
	return updateToken(t.token)
}

func (t Token) validateRequest(w *http.ResponseWriter, rotate bool) bool {
	token := Token{t.token}
	authed, err := token.verify()
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return false
	} else if !authed {
		(*w).WriteHeader(http.StatusUnauthorized)
		return false
	} else if rotate {
		if err = token.updateToken(); err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			return false
		}
	}
	return true
}
