package handlefuncs

import (
	"encoding/json"
	"net/http"
	dbfuncs "server/pkg/db/dbfuncs"
)

type DataID struct {
	Id string `json:"id"`
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

    users, err dbfuncs.Getusers()


	} else {

		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}

}