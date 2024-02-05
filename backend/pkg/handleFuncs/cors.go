package handlefuncs

import (
	"fmt"
	"net/http"
)

func Cors(w *http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	fmt.Println(origin, "cors line 10")
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, user_token")

	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == http.MethodOptions {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}
