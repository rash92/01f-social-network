package handlefuncs

import (
	"encoding/json"
	"net/http"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		posts, err := FindAllChats()

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(posts)

	} else {

		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
