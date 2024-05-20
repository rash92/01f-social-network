package handlefuncs

import (
	"encoding/json"
	"net/http"
)

func HandleCatogries(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)


	if r.Method == http.MethodGet {

		Categories, err := FindAllCats()

		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(Categories)

	} else {

		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
