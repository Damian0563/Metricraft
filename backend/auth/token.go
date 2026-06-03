package auth

import (
	"backend/db"
	"net/http"
	"strings"
)

type Token struct {
	token string
}

func NewToken(token string) Token {
	return Token{token: token}
}

func (t Token) GetAppName() (string, error) {
	split := strings.Split(t.token, ":")
	return db.GetAppNameByToken(split[0])
}

func (t Token) Sign(rotate bool) (string, error) {
	split := strings.Split(t.token, ":")
	return db.SignToken(split[0], rotate)
}

func (t Token) Verify() (bool, error) {
	return db.CheckToken(t.token)
}

func (t Token) UpdateToken() error {
	return db.UpdateToken(t.token)
}

func (t Token) ValidateRequest(w *http.ResponseWriter, rotate bool) bool {
	token := NewToken(t.token)
	authed, err := token.Verify()
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return false
	} else if !authed {
		(*w).WriteHeader(http.StatusUnauthorized)
		return false
	} else if rotate {
		if err = token.UpdateToken(); err != nil {
			(*w).WriteHeader(http.StatusInternalServerError)
			return false
		}
	}
	return true
}
