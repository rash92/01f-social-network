package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/pkg/db/dbfuncs"
)

func HandleValidateCookie(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("user_token")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}
		valid, err := dbfuncs.ValidateCookie(cookie.Value)
		if !valid || err != nil {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		var username string
		var imgURL string
		var id string

		// var storedPassword string

		query := `
    SELECT Sessions.user, Sessions.userId, Users.profileImg
FROM Sessions
JOIN Users ON Sessions.UserID = Users.Id
WHERE Sessions.Id = ?

`
		err = db.QueryRow(query, cookie.Value).Scan(&username, &id, &imgURL)

		if err != nil {
			http.Error(w, `{"error": "something went wrong "}`, http.StatusBadRequest)
			return
		}

		response := map[string]interface{}{
			"success":    true,
			"username":   username,
			"profileImg": imgURL,
			"id":         id,
		}
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)

	}
}
