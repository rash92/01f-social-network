package handlefuncs

import (
	"encoding/json"
	"fmt"
	"net/http"
	dbfuncs "server/pkg/db/dbfuncs"
)

func HandleValidateCookie(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("user_token")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}
		if !dbfuncs.ValidateCookie(cookie.Value) {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		var username string
		var imgURL string
		var id string

		// var storedPassword string

		query := `
    SELECT  Sessions.userId, Users.Profile
FROM Sessions
JOIN Users ON Sessions.UserID = Users.Id
WHERE Sessions.Id = ?

`
		err = database.QueryRow(query, cookie.Value).Scan( &id, &imgURL)

		if err != nil {
			fmt.Println(err, "\nFailed to get user info for validating token.")
			http.Error(w, `{"error": "something went wrong "}`, http.StatusBadRequest)
			return
		}

		
		
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)

	}
}