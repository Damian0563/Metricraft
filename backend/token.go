package main

func (t Token) GetUser() (User, error) {
	return getUserByToken(t.token)
}

func (t Token) GetAppName() (string, error) {
	return getAppNameByToken(t.token)
}

func (t Token) ChangeRealtime(enabled bool) error {
	return changeRealtimeByToken(t.token, enabled)
}

func (t Token) sign() (string, error) {
	return signToken(t.token)
}

func (t Token) verify() (bool, error) {
	return checkToken(t.token)
}

func (t Token) updateToken(signed string) error {
	return updateToken(t.token, signed)
}
