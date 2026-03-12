package enter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Enter(w http.ResponseWriter, r *http.Request) {
	//inspect the request
	var headers map[string]any = make(map[string]any)
	headers["X-Forwarded-For"] = r.Header.Get("X-Forwarded-For")
	headers["X-Forwarded-Host"] = r.Header.Get("X-Forwarded-Host")
	headers["X-Forwarded-Proto"] = r.Header.Get("X-Forwarded-Proto")
	headers["X-Real-IP"] = r.Header.Get("X-Real-IP")
	method := r.Method
	var data []byte
	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&data)
	if err == nil {
		json.Unmarshal(data, &body)
	}
	payload := Payload{headers, body, method}
	fmt.Println(payload)
	redirect := strings.Replace(r.URL.String(), os.Getenv("PORT"), os.Getenv("DEST_PORT"), 1)
	fmt.Println(redirect)
	//switch method {
	// case "GET":
	// 	http.Get(r)
	// case "POST":
	// 	Post()
	// case "PUT":
	// 	Put()
	// case "DELETE":
	// 	Delete()
	// }
	Leave(payload) //send data to frontend
}
