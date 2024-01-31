package handlefuncs

import (
	"encoding/json"
	"net/http"
	"server/pkg/db/dbfuncs"
)

type reaction struct {
	Postid string `json:"postId"`
	Query  string `json:"query"`
}

func HandlePostLikeDislike(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	if r.Method == http.MethodPost {
		var entredData reaction
		errj := json.NewDecoder(r.Body).Decode(&entredData)
		if errj != nil {
			http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("user_token")
		if err != nil {
			http.Error(w, `{"error": "something went wrong please login"}`, http.StatusUnauthorized)
			return
		}
		if !dbfuncs.ValidateCookie(cookie.Value) {
			http.Error(w, `{"error": "something went wrong please login"}`, http.StatusUnauthorized)
			return
		}

		var userID string

		err = database.QueryRow("SELECT UserId FROM Sessions WHERE Id=?", cookie.Value).Scan(&userID)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		if entredData.Query == "like" {
			err := dbfuncs.AddLikes(userID, entredData.Postid)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			like, dislikes := dbfuncs.CountLikesDislikes(entredData.Postid)
			userlikes := dbfuncs.FindLikeUsers(entredData.Postid)
			response := map[string]interface{}{
				"likes":     like,
				"dislikes":  dislikes,
				"userlikes": userlikes,
				"id":        entredData.Postid,
			}

			json.NewEncoder(w).Encode(response)

		} else {
			err := dbfuncs.AddDislikes(userID, entredData.Postid)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			like, dislikes := dbfuncs.CountLikesDislikes(entredData.Postid)
			userlikes := dbfuncs.FindLikeUsers(entredData.Postid)
			response := map[string]interface{}{
				"likes":     like,
				"dislikes":  dislikes,
				"userlikes": userlikes,
				"id":        entredData.Postid,
			}
			json.NewEncoder(w).Encode(response)

		}

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}
