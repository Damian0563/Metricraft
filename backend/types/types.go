package types

type PendingUsers struct {
	Mail string `json:"mail"`
}

type AllowedUsersDb struct {
	Mail                 string `json:"mail"`
	ReceiveNotifications bool   `json:"receiveNotifications"`
}

type AllowedUsers struct {
	Mail                 string `json:"mail"`
	Initials             string `json:"initials"`
	Status               bool   `json:"status"`
	ReceiveNotifications bool   `json:"receiveNotifications"`
}
type SendRecoveryUser struct {
	Mail string `json:"mail"`
	Id   string `json:"id,omitempty"`
}

type RecoveryCheckPayload struct {
	Password string `json:"password"`
}
type Invite struct {
	Invitees []string `json:"invitees"`
}

type ExistsErrResponse struct {
	Exists bool
	Err    error
	Origin string
	Owner  string
}
