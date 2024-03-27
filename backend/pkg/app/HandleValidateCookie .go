package app

// import (
// 	"backend/pkg/db/dbfuncs"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// func HandleValidateCookie(w http.ResponseWriter, r *http.Request) {
// 	Cors(&w, r)
// 	// if r.Method == http.MethodOptions {
// 	// 	w.WriteHeader(http.StatusOK)
// 	// 	return
// 	// }
// 	if r.Method == http.MethodGet {
// 		cookie, err := r.Cookie("user_token")
// 		if err != nil {
// 			fmt.Println(err)
// 			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		valid, err := dbfuncs.ValidateCookie(cookie.Value)
// 		if !valid || err != nil {
// 			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// var storedPassword string

// }
// }
