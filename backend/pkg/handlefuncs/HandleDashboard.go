package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	var ownerId string

	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&ownerId)
	if err != nil {
		http.Error(w, "something went wrong ", http.StatusInternalServerError)
		return
	}

	chatIDs, err := GetAllChats(ownerId)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var users []dbfuncs.User

	for _, chat := range chatIDs {
		fmt.Println("chat", chat, "chatIDs", chatIDs)

		if chat != "" {
			user, err := dbfuncs.GetUserById(chat)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				log.Println("error getting user from id", err)
				return
			}
			users = append(users, user)
		}

	}

	group, err := GetGroupDash(ownerId)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	posts, err := dbfuncs.GetVisiblePosts(ownerId)
	if err != nil {
		log.Println("error getting posts :)", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	nofications, err := dbfuncs.GetAllNotificationsByRecieverId(ownerId)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"chats":         users,
		"groups":        group.GroupCards,
		"notifications": nofications,
		"Posts":         posts,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("error sending chats")
		return
	}
}
