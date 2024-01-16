package handlefuncs

import (
	"encoding/json"
	"fmt"
	"net/http"
	dbfuncs "server/pkg/db/dbfuncs"
)

func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)

	if r.Method == http.MethodPost {
		var newComment Comment
		errj := json.NewDecoder(r.Body).Decode(&newComment)
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
		if len(newComment.Body) > CharacterLimit {

			http.Error(w, `{"error": "413 Payload Too Large"}`, http.StatusRequestEntityTooLarge)
			return
		}
		if len(newComment.Body) == 0 {

			http.Error(w, `{"error": "204 No Content"}`, http.StatusNoContent)
			return
		}

		PostId := newComment.PostID
		PostBody := newComment.Body

		_, err = dbfuncs.AddComment(PostBody, cookie.Value, PostId)
		// fmt.Println((commentId))
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		comment, err := FindPostsComment(PostId)

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		fmt.Println(comment)
		response := map[string]interface{}{
			"success":  true,
			"comments": comment,
			"id":       PostId,
		}
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}
