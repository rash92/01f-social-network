package handlefuncs

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type data struct {
	Page      int    `json:"page"`
	CurrUser  string `json:"currUser"`
	OtherUser string `json:"otherUser"`
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	Cors(&w, r)
	if r.Method == http.MethodPost {
		var entredData data
		errj := json.NewDecoder(r.Body).Decode(&entredData)
		if errj != nil {
			http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
			return
		}
		// page, _ := strconv.Atoi(entredData.Page)
		batchSize := 10

		offset := (entredData.Page - 1) * batchSize

		db, err := sql.Open("sqlite3", "../sever/forum.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		query := `
    SELECT * FROM Messages
    WHERE (SenderId = ? AND RecipientId = ?) OR (SenderId = ? AND RecipientId = ?)
    ORDER BY Created DESC
    LIMIT ? OFFSET ?
`
		rows, err := database.Query(query, entredData.CurrUser, entredData.OtherUser, entredData.OtherUser, entredData.CurrUser, batchSize, offset)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		messages := []Message{}
		for rows.Next() {
			var message Message
			err := rows.Scan(&message.ID, &message.SenderID, &message.RecipientID, &message.Message, &message.Created)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			messages = append(messages, message)
		}

		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}

		response := struct {
			Messages []Message `json:"messages"`
		}{
			Messages: messages,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
