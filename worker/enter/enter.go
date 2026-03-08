package enter

import (
	"fmt"
	"net/http"
)

func Enter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
