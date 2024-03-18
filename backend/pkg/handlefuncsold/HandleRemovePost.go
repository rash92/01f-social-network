package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"net/http"
)

type postId struct {
	PostId string `json:"postId"`
}

func HandleRemovePost(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	if r.Method == http.MethodPost {

		var entredData postId
		errj := json.NewDecoder(r.Body).Decode(&entredData)
		if errj != nil {
			http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("user_token")
		if err != nil {
			fmt.Println(err)
			http.Error(w, `{"error": "something went wrong please login"}`, http.StatusUnauthorized)
			return
		}
		if !dbfuncs.ValidateCookie(cookie.Value) {
			http.Error(w, `{"error": "something went wrong please login"}`, http.StatusUnauthorized)
			return
		}

		userIdSession, err := dbfuncs.GetUserIdFromCookie(cookie.Value)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		userIdPost, err := dbfuncs.GetUserIdFromPostId(entredData.PostId)

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		if userIdSession == userIdPost {
			dbfuncs.RemovePost(entredData.PostId)

			response := map[string]interface{}{
				"success": true,
				"postId":  entredData.PostId,
			}
			json.NewEncoder(w).Encode(response)

		}

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

}
