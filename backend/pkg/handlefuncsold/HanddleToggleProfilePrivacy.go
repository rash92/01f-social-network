package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"net/http"
)

type PrivcySetting struct {
	Privacy string `json:"setting"`
	UserId  string `json:"id"`
}

func HanddleToggleProfilePrivacy(w http.ResponseWriter, r *http.Request) {

	var privacy PrivcySetting

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&privacy)

	if err != nil {
		errorMessage := fmt.Sprintf("error decoding userId: %v", err.Error())
		fmt.Println(err.Error(), "60")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	fmt.Println("privacy setting", privacy)
	err = dbfuncs.UpdatePrivacySetting(privacy.UserId, privacy.Privacy)
	if err != nil {
		errorMessage := fmt.Sprintf("error updating privacy setting: %v", err.Error())
		fmt.Println(err.Error(), "66")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	response := make(map[string]string)
	response["message"] = "Successfully updated privacy setting"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
