package handlefuncs

import (
	"encoding/json"
	"net/http"
	dbfuncs "backend/pkg/db/dbfuncs"
)

type DataID struct {
	Id string `json:"id"`
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		users, err := dbfuncs.Getusers()
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(users)

	} else {

		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
