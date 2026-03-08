package main

import (
	"github.com/joho/godotenv"
	"metricraft/worker/enter"
	"metricraft/worker/leave"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	if !(os.Getenv("PORT") == "" || os.Getenv("SECRET") == "") {
		panic("Port and secret must be set in .env file")
	}
	http.HandleFunc("/enter", enter.Enter)
	leave.Leave()
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
