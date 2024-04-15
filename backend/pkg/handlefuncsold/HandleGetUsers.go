package handlefuncs

import (
	dbfuncs "backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"net/http"
)

type DataID struct {
	Id string `json:"id"`
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		users, err := dbfuncs.GetUsers()
		fmt.Println(users, "users")
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(users)

	} else {

		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
