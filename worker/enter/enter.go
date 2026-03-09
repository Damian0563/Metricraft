package enter

import (
	"fmt"
	"metricraft/worker/leave"
	"net/http"
)

func Enter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
	//inspect the request
	leave.Leave() //send data to frontend
	//reroute request to dest port
}
