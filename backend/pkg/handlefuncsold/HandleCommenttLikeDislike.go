package handlefuncs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/pkg/db/dbfuncs"
)

type Commentreaction struct {
	Postid    string `json:"postId"`
	CommentId string `json:"commentId"`
	Query     string `json:"query"`
}

func HandleCommenttLikeDislike(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	if r.Method == http.MethodPost {
		var entredData Commentreaction

		errj := json.NewDecoder(r.Body).Decode(&entredData)

		fmt.Println(entredData, "checking comment")
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
			dbfuncs.LikeComment(userID, entredData.CommentId)
		} else {
			dbfuncs.DislikeComment(userID, entredData.CommentId)
		}

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		like, dislikes := dbfuncs.CountCommentReacts(entredData.CommentId)
		response := map[string]interface{}{
			"likes":    like,
			"dislikes": dislikes,
		}

		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}
