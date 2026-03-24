package backend

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	router := http.NewServeMux()
	port := os.Getenv("PORT")
	router.HandleFunc("/", welcome)

	http.ListenAndServe(":"+port, router)
	fmt.Println("Server started")
}
