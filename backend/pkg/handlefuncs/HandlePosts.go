package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"net/http"
)

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
	posts, err := dbfuncs.GetPostById(postQuery.UserId, postQuery.PostID)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(posts)

}

func HandlePostLikeDislike(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return

	}
	var enteredData Reaction
	errj := json.NewDecoder(r.Body).Decode(&enteredData)
	if errj != nil {
		http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
		return
	}

	err := dbfuncs.LikeDislikePost(enteredData.UserId, enteredData.Postid, enteredData.Query)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return

	}

	like, dislikes, err := dbfuncs.CountPostReacts(enteredData.Postid)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return

	}

	UserLikeDislike, err := dbfuncs.GetUserLikeDislike(enteredData.UserId, enteredData.Postid)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return

	}

	response := map[string]interface{}{
		"Likes":           like,
		"Dislikes":        dislikes,
		"UserLikeDislike": UserLikeDislike,
		"id":              enteredData.Postid,
	}
	json.NewEncoder(w).Encode(response)

}
