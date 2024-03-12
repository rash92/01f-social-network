package handlefuncs

import (
	"encoding/json"
	"net/http"
	"server/pkg/db/dbfuncs"
)





func GetAllChats(ownerId) ([]string, err) {
    var chatIds []string

    messaged, err := dbfuncs.GetAllUserIdsSortedByLastPrivateMessage()
    if err != nil {
        return chatIds, err
    }

    unmessaged, err = dbfuncs.GetUnmessagedUserIdsSortedAlphabetically()
    if err != nil {
        return chatIds, err
    }
    
    chats = append(messaged, unmessaged...)
    return chatIds, nil
}

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
    var ownerId string

    if r.Method != http.MethodPost {
        http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    err := json.NewDecoder(r.Body).Decode(&ownerId)
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

    chatIDs, err := GetAllChats(ownerId)
	if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		}

    var users []dbfuncs.User
    for _, chat := range chatIDs {
        user, err := dbfuncs.GetUserById()
        if err != nil {
            log.Println("error getting user from id")
            return
        }
        users = append(users, user)
    }

    response := map[string]interface{}{
		"chats": users,
	}

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Println("error sending chats")
            return
    }
}