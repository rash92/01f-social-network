package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"net/http"
)

type searchFollowQuery struct {
	Search string `json:"search"`
	Id     string `json:"id"`
}

func HandleSearchFollower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var searchQuery searchFollowQuery

	err := json.NewDecoder(r.Body).Decode(&searchQuery)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := dbfuncs.SearchFollowers(searchQuery.Search, searchQuery.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error": "Error getting user from database"}`, http.StatusInternalServerError)
		return
	}

	// Send the user back to the client
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, `{"error": "Error encoding JSON"}`, http.StatusInternalServerError)
		return
	}
}
