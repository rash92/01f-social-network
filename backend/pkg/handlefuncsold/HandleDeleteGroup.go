package handlefuncs

import (
	"encoding/json"
	"net/http"
)

type Deletgroup struct {
	Groupid string `json:"groupid"`
	UserId  string `json:"userid"`
}

func HandleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var deleteType Deletgroup
	err := json.NewDecoder(r.Body).Decode(&deleteType)
	if err != nil {

		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	group, err := GetGroup(deleteType.Groupid, deleteType.UserId)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
	}

	if group.BasicInfo.CreatorId != deleteType.UserId {
		http.Error(w, `{"error": "you not creator of this group"}`, http.StatusBadRequest)
	}
	response := map[string]interface{}{
		"success": true,
	}

	json.NewEncoder(w).Encode(response)

	w.WriteHeader(http.StatusOK)

}
