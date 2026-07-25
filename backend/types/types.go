package types

type DashboardInitPayload struct {
	AppName      string   `json:"appName"`
	SignedSecret string   `json:"signedSecret"`
	Settings     Settings `json:"settings"`
	Error        string   `json:"error"`
	Urls         []string `json:"urls"`
}

type Settings struct {
	Realtime  bool                     `json:"realtime"`
	Enabled   map[string]EnabledMetric `json:"enabled"`
	Retention int                      `json:"retention"`
}

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
