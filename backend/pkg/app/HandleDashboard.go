package app

import (
    "encoding/json"
    "log"
    "net/http"
    "server/pkg/db/dbfuncs"
)

func GetAllChats(ownerId string ) ([]string, error) {
    var chatIds []string

    messaged, err := dbfuncs.GetAllUserIdsSortedByLastPrivateMessage(ownerId)
    if err != nil {
        return chatIds, err
    }

    unmessaged, err := dbfuncs.GetUnmessagedUserIdsSortedAlphabetically(ownerId)
    if err != nil {
        return chatIds, err
    }

    chats := append(messaged, unmessaged...)
    return chats, nil
}

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
    }

    var users []dbfuncs.User
    for _, chat := range chatIDs {
        user, err := dbfuncs.GetUserById(chat)
        if err != nil {
            log.Println("error getting user from id")
            return
        }
        users = append(users, user)
    }

	// get all posts you're allowed to see and all groups.

    response := map[string]interface{}{
        "chats": users,
    }

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Println("error sending chats")
        return
    }
}