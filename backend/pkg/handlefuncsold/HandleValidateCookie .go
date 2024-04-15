package handlefuncs

import (
	dbfuncs "backend/pkg/db/dbfuncs"
	"encoding/json"
	"net/http"
)

func HandleValidateCookie(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == http.MethodGet {
		cookie, err:= r.Cookie("user_token")

		// 		var username string
		// 		var imgURL string
		// 		var id string

		// 		// var storedPassword string

		// 		query := `
		//     SELECT  Sessions.userId, Users.Avatar
		// FROM Sessions
		// JOIN Users ON Sessions.UserID = Users.Id
		// WHERE Sessions.Id = ?

		// `
		// 		err = database.QueryRow(query, cookie.Value).Scan(&id, &imgURL)

		//	if err != nil {
		//		fmt.Println(err, "\nFailed to get user info for validating token.")
		//		http.Error(w, `{"error": "something went wrong "}`, http.StatusBadRequest)
		//		return
		//	}



		if cookie == nil || err != nil{
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return

		}
		id, err := dbfuncs.GetUserIdFromCookie(cookie.Value)

		if err != nil {

			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		user, err := dbfuncs.GetUserById(id)

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"user":    user,
		}

		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)

	}
}
