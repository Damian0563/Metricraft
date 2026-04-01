package main

import "fmt"

func createUser(mail string, secret string, appName string) {
	fmt.Printf(mail, secret, appName)
}

func signIn(mail string, secret string) {
	fmt.Printf(mail, secret)
}
