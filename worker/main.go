package main

import (
	"github.com/joho/godotenv"
	"metricraft/worker/enter"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	if os.Getenv("PORT") == "" || os.Getenv("SECRET") == "" || os.Getenv("DEST_PORT") == "" {
		panic("Port and secret must be set")
	}
	http.HandleFunc("/", enter.Enter)
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
