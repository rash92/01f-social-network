package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"net/http"
	"strings"
)

func HanndleUserNameIsDbOrEmail(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	if r.Method == "POST" {
		val := strings.Replace(r.URL.Path, "/check-", "", 1)

		username := r.FormValue(val)

		msg, isexist, err := dbfuncs.CheckValueInDB(w, r, username, val)

		if err != nil {
			http.Error(w, msg, http.StatusInternalServerError)

			return
		}
		// // Return the result to the frontend
		response := map[string]bool{"exists": isexist}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
