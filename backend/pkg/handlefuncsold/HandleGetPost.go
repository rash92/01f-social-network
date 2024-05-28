package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"net/http"
)

type PostQuery struct {
	PostID string `json:"post_id"`
	UserId string `json:"user_id"`
}

func HandleGetPost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var postQuery PostQuery
	err := json.NewDecoder(r.Body).Decode(&postQuery)
	if err != nil {

		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	posts, err := dbfuncs.GetPostById(postQuery.UserId,  postQuery.PostID)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(posts)

}
