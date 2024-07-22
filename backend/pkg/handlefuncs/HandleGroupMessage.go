package handlefuncs

// import (
// 	"backend/pkg/db/dbfuncs"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// func GroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var enteredData MessageData
// 	err := json.NewDecoder(r.Body).Decode(&enteredData)
// 	if err != nil {
// 		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
// 		return
// 	}
// 	var messages []dbfuncs.G
// 	if enteredData.Type == "GroupMessage" {

// 		messages, err = dbfuncs.GetAllPrivateMessagesByUserId(enteredData.CurrUser, enteredData.OtherUser)
// 		if err != nil {
// 			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 			return
// 		}

// 		// else {
// 		// 	// Why is this branch empty? Was it empty when I downloaded the repo
// 		// 	// but just didn't register as a problem till I saved some comments
// 		// 	// in a different file? It's not necessary, anyway, after the return.
// 		// }

// 	}

// 	response := map[string]interface{}{
// 		"success":  true,
// 		"messages": messages,
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)

// }
