package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"net/http"
)

type data struct {
	CurrUser  string `json:"currUser"`
	OtherUser string `json:"otherUser"`
	Type      string `json:"type"`
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var entredData data
	err := json.NewDecoder(r.Body).Decode(&entredData)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	var messages []dbfuncs.PrivateMessage
	if entredData.Type == "privateMessage" {
		fmt.Println(entredData.Type, "entredData.Type")

		messages, err = dbfuncs.GetAllPrivateMessagesByUserId(entredData.CurrUser, entredData.OtherUser)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

	}

	response := map[string]interface{}{
		"success":  true,
		"messages": messages,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
