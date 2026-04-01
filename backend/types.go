package main

type signPayload struct {
	Mail    string `json:"mail"`
	Secret  string `json:"secret"`
	AppName string `json:"appName",omitempty`
}
