package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"net/http"
)

// func MessagesHandlerNew(w http.ResponseWriter, r *http.Request) {
// 	Cors(&w, r)
// 	if r.Method == http.MethodPost {
// 		var enteredData MessageData
// 		err := json.NewDecoder(r.Body).Decode(&enteredData)
// 		if err != nil {
// 			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
// 			return
// 		}
// 		batchSize := 10

// 		offset := (enteredData.Page - 1) * batchSize

// 		dbMessages, err := dbfuncs.GetLimitedPrivateMessages(enteredData.CurrUser, enteredData.OtherUser, batchSize, offset)
// 		if err != nil {
// 			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 			return
// 		}
// 		var frontEndMessages []PrivateMessage
// 		for _, dbMessage := range dbMessages {
// 			frontendMessage := PrivateMessage{
// 				Id:          dbMessage.Id,
// 				SenderId:    dbMessage.SenderId,
// 				RecipientId: dbMessage.RecipientId,
// 				Message:     dbMessage.Message,
// 				CreatedAt:   dbMessage.CreatedAt,
// 			}
// 			frontEndMessages = append(frontEndMessages, frontendMessage)
// 		}

// 		response := struct {
// 			Messages []PrivateMessage `json:"messages"`
// 		}{
// 			Messages: frontEndMessages,
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)

// 	}
// }

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var enteredData MessageData
	err := json.NewDecoder(r.Body).Decode(&enteredData)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	var messages []dbfuncs.PrivateMessage
	if enteredData.Type == "privateMessage" {
		fmt.Println(enteredData.Type, "entredData.Type")

		messages, err = dbfuncs.GetAllPrivateMessagesByUserId(enteredData.CurrUser, enteredData.OtherUser)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}else{
			
		}

	}

	response := map[string]interface{}{
		"success":  true,
		"messages": messages,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// func MessagesHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("started messages handler")
// 	Cors(&w, r)
// 	fmt.Println("started messages handler after cors")
// 	if r.Method == http.MethodPost {
// 		fmt.Println("started messages handler in method == post")
// 		var enteredData MessageData
// 		errj := json.NewDecoder(r.Body).Decode(&enteredData)
// 		if errj != nil {
// 			http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
// 			return
// 		}
// 		batchSize := 10

// 		offset := (enteredData.Page - 1) * batchSize
// 		fmt.Println("started messages handler before opening fake database")
// 		db, err := sql.Open("sqlite3", "../sever/forum.db")
// 		fmt.Println("started messages handler after opening fake database: ", err)
// 		if err != nil {

// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			fmt.Println("internal server error from opening fake database")
// 			return
// 		}
// 		defer db.Close()

// 		query := `
// 				SELECT * FROM Messages
// 				WHERE (SenderId = ? AND RecipientId = ?) OR (SenderId = ? AND RecipientId = ?)
// 				ORDER BY Created DESC
// 				LIMIT ? OFFSET ?
// 			`
// 		rows, err := db.Query(query, enteredData.CurrUser, enteredData.OtherUser, enteredData.OtherUser, enteredData.CurrUser, batchSize, offset)
// 		if err != nil {
// 			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
// 			fmt.Println("internal server error from querying fake database so select messages")
// 			return
// 		}
// 		defer rows.Close()

// 		messages := []Message{}
// 		for rows.Next() {
// 			var message Message
// 			err := rows.Scan(&message.ID, &message.SenderID, &message.RecipientID, &message.Message, &message.Created)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			messages = append(messages, message)
// 		}

// 		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
// 			messages[i], messages[j] = messages[j], messages[i]
// 		}

// 		response := struct {
// 			Messages []Message `json:"messages"`
// 		}{
// 			Messages: messages,
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 	}
// }

func GetAllChats(ownerId string) ([]string, error) {

	var chatIds []string

	messaged, err := dbfuncs.GetAllUserIdsSortedByLastPrivateMessage(ownerId)
	if err != nil {
		fmt.Println("error getting all user ids sorted by last private message", err)
		return chatIds, err
	}

	return messaged, err
}
