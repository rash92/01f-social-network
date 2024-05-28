package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"net/http"
)

type reaction struct {
	Postid string `json:"postId"`
	Query  string `json:"query"`
	UserId string `json:"id"`
}

func HandlePostLikeDislike(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return

	}
	var entredData reaction
	errj := json.NewDecoder(r.Body).Decode(&entredData)
	if errj != nil {
		http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
		return
	}

	err := dbfuncs.LikeDislikePost(entredData.UserId, entredData.Postid, entredData.Query)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return

	}

	like, dislikes, err := dbfuncs.CountPostReacts(entredData.Postid)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return

	}

	UserLikeDislike, err := dbfuncs.GetUserLikeDislike(entredData.UserId, entredData.Postid)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return

	}

	response := map[string]interface{}{
		"Likes":           like,
		"Dislikes":        dislikes,
		"UserLikeDislike": UserLikeDislike,
		"id":              entredData.Postid,
	}
	json.NewEncoder(w).Encode(response)

}
