package main

import "strings"

func (t Token) GetUser() (User, error) {
	return getUserByToken(t.token)
}

func (t Token) GetAppName() (string, error) {
	split := strings.Split(t.token, ":")
	return getAppNameByToken(split[0])
}

func (t Token) ChangeRealtime(enabled bool) error {
	split := strings.Split(t.token, ":")
	return changeRealtimeByToken(split[0], enabled)
}

func (t Token) sign() (string, error) {
	split := strings.Split(t.token, ":")
	return signToken(split[0])
}

func (t Token) verify() (bool, error) {
	return checkToken(t.token)
}

func (t Token) updateToken(signed string) error {
	return updateToken(t.token, signed)
}
