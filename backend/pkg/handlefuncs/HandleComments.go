package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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

type file struct {
	bytes     []byte
	extension string
}

type validatedCommentRequest struct {
	Comment
	image *file
}

func validateCommentRequest(r *http.Request) (*validatedCommentRequest, error) {

	if r.Method != http.MethodPost {
		return nil, errors.New(`{"error": "405 Method Not Allowed"}`)
	}

	validated := validatedCommentRequest{}
	err := json.NewDecoder(r.Body).Decode(&validated)
	if err != nil {
		return nil, errors.New(`{"error": "` + err.Error() + `"}`)
	}

	if len(validated.Body) > CharacterLimit {
		return nil, errors.New(`{"error": "413 Payload Too Large"}`)
	}

	if len(validated.Body) == 0 {

		return nil, errors.New(`{"error": "204 No Content"}`)
	}

	if validated.Image != "" {
		imageFile, err := ConvertBase64ToImage(validated.Image)
		if err != nil {
			return nil, err
		}
		validated.image = imageFile
	}

	return &validated, nil
}

func ConvertBase64ToImage(base64String string) (*file, error) {
	// Split the base64 string to isolate the MIME type and the actual data

	splitData := strings.Split(base64String, ",")
	if len(splitData) != 2 {
		return nil, fmt.Errorf("invalid base64 string")
	}

	mimeType := strings.Split(splitData[0], ";")[0]
	data := splitData[1]

	// Map the MIME type to a file extension
	mimeToExtension := map[string]string{
		"data:image/jpeg": ".jpg",
		"data:image/png":  ".png",
		"data:image/gif":  ".gif",
		// Add more mappings as needed
	}
	extension, ok := mimeToExtension[mimeType]
	if !ok {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Decode the base64 string back to bytes
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	image := file{
		decodedData,
		extension,
	}

	return &image, nil
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
