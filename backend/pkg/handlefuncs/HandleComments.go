package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func DbMessagesToFrontend(dbMessages []dbfuncs.GroupMessage) []GroupMessage {
	var frontendGroupMessages []GroupMessage
	for _, dbGroupMessage := range dbMessages {
		frontendGroupMessages = append(frontendGroupMessages, DbGroupMessageToFrontend(dbGroupMessage))
	}
	return frontendGroupMessages
}

func DbMessageToFrontend(dbMessage dbfuncs.GroupMessage) GroupMessage {
	frontendGroupMessage := GroupMessage{
		Id:        dbMessage.Id,
		SenderId:  dbMessage.SenderId,
		GroupId:   dbMessage.GroupId,
		Message:   dbMessage.Message,
		CreatedAt: dbMessage.CreatedAt,
	}
	return frontendGroupMessage
}

func HandleAddComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "405 Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return

	}
	var newComment Comment
	errj := json.NewDecoder(r.Body).Decode(&newComment)
	if errj != nil {
		http.Error(w, `{"error": "`+errj.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if len(newComment.Body) > CharacterLimit {
		http.Error(w, `{"error": "413 Payload Too Large"}`, http.StatusRequestEntityTooLarge)
		return
	}
	if len(newComment.Body) == 0 {

		http.Error(w, `{"error": "204 No Content"}`, http.StatusNoContent)
		return
	}

	newCommentDb := dbfuncs.Comment{
		Body:            newComment.Body,
		CreatorId:       newComment.CreatorId,
		PostId:          newComment.PostID,
		CreatedAt:       time.Now(),
		Likes:           0,
		Dislikes:        0,
		CreatorNickname: newComment.CreatorNickname,
	}

	if newComment.Image != "" {
		imageUUID, err := dbfuncs.ConvertBase64ToImage(newComment.Image, "./pkg/db/images")
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			log.Println("error converting base64 to image", err)
			return
		}
		newCommentDb.Image = imageUUID
	}

	id, err := dbfuncs.AddComment(&newCommentDb)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	newCommentDb.Id = id

	response := map[string]interface{}{
		"success": true,
		"comment": newCommentDb,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
