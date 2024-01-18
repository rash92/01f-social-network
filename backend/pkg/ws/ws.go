package ws

//general notes/ ideas:
//consider redoing database tables to combine private messages, group messages and various types of notificaitons in one table
//and differentiate them with a type field, possibly with empty fields for fields that are not needed for that type of message

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"server/pkg/db/dbfuncs"
	"server/pkg/handlefuncs"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var db *sql.DB

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:8000"

		},
	}

	//need to flip this map to be map[user]*websocket.Conn (or userid or string etc.)
	activeConnections = make(map[*websocket.Conn]string)
	connectionLock    sync.Mutex
	dbLock            sync.Mutex

	// userListLock      sync.Mutex
)

//check capitalization of json bit

// this is what the client sends to the server,
// Body will be reunmarshalled based on type into PrivateMessage, GroupMessage, or Notification etc.
// as well as ws messages unrelated to database operations
type WsMessage struct {
	Type      string    `json:"type"`
	Body      []byte    `json:"message"`
	TimeStamp time.Time `json:"time"`
}

// these should match the database fields, check types in database and make them match, also move them to db package eventually
type PrivateMessage struct {
	Id          string    `json:"Id"`
	SenderId    string    `json:"SenderId"`
	RecipientId string    `json:"RecipientId"`
	Message     string    `json:"Message"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type GroupMessage struct {
	Id        string    `json:"Id"`
	SenderId  string    `json:"SenderId"`
	GroupId   string    `json:"GroupId"`
	Message   string    `json:"Message"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type Notification struct {
	Id         string    `json:"Id"`
	RecieverId string    `json:"RecieverId"`
	SenderId   string    `json:"SenderId"`
	Body       string    `json:"Body"`
	Type       string    `json:"Type"`
	Seen       bool      `json:"Seen"`
	CreatedAt  time.Time `json:"CreatedAt"`
}

type Post struct {
	//fill in later
}

type Comment struct {
	//fill in later
}

//these may not involve database calls but can still be sent through websockets

// this can be resused for sending WsMessages to other users about a user who has
// e.g. registered, logged in, logged out, or changed their status
type BasicUserInfo struct {
	UserId         string `json:"UserId"`
	PrivacySetting string `json:"PrivacySetting"` //maybe?
	//fill in later if more needed
}

// pull some of this stuff out into separate functions to make cleaner
// can consider channel approach instead of mutex
func handleConnection(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_token")
	if err != nil || !dbfuncs.ValidateCookie(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	userID, err := dbfuncs.GetUserIdFromCokie(cookie.Value) // fix spelling of cokie
	if err != nil {
		log.Println("Error retrieving userID from database:", err)
	}

	connectionLock.Lock()
	activeConnections[conn] = userID
	connectionLock.Unlock()

	for {
		_, msgBytes, err := conn.ReadMessage()
		//possibly don't want to immediately delete the connection if there is an error
		if err != nil {
			connectionLock.Lock()
			delete(activeConnections, conn)
			connectionLock.Unlock()
			log.Println("User", userID, "disconnected, unable to read from websocket, error:", err)
			break
		}
		var recievedData WsMessage
		err = json.Unmarshal(msgBytes, &recievedData)
		if err != nil {
			log.Println("Error unmarshalling websocket message:", err)
		}

		switch recievedData.Type {
		case "privateMessage":
			//fill in
		case "groupMessage":
			//fill in
		case "notification":
			//fill in
		case "post":
			//fill in
		case "comment":
			//fill in
		case "logout":
			//fill in
		case "login":
			//fill in
		case "register":
			//fill in
		case "updatePrivacySetting":
			//fill in
		default:
			//unexpected type
			log.Println("Unexpected websocket message type:", recievedData.Type)
		}

	}

}

func handleConnectionOld(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_token")
	if err != nil || !dbfuncs.ValidateCookie(cookie.Value) {

		// If the cookie is not valid, close the WebSocket connection
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer func() {
		// fmt.Println("conncetion closing")
		conn.Close()

	}()

	userID, err := dbfuncs.GetUserIdFromCokie(cookie.Value)
	if err != nil {
		fmt.Println("Error sending user ID to client:", err)
	}

	connectionLock.Lock()
	activeConnections[conn] = userID
	connectionLock.Unlock()

	isLoggingOut := false
	broadcastUserList()

	for {

		_, p, err := conn.ReadMessage()
		if err != nil {
			connectionLock.Lock()
			if isLoggingOut {
				for client, id := range activeConnections {
					if id == userID && client != conn {
						fmt.Println(id, client)
						data := map[string]interface{}{
							"data": "log me out",
							"type": "logout",
						}
						err := client.WriteJSON(data)
						if err != nil {
							fmt.Println("Error logout message  to client:", err)
						}
					}
				}

			}
			delete(activeConnections, conn)
			connectionLock.Unlock()
			broadcastUserList()
			fmt.Println("User", userID, "disconnected")

			break
		}

		var receivedData handlefuncs.Message

		err = json.Unmarshal(p, &receivedData)
		if receivedData.Type == "logout" {
			isLoggingOut = true
		}
		if err != nil {
			fmt.Println(err, "eror of  Unmarshal")
		}
		fmt.Println(receivedData, "receivedData message")
		dbLock.Lock()
		id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
		dbLock.Unlock()
		if err != nil {
			fmt.Println(err, "error adding message in the data base main line 120")
		}

		receivedData.ID = id.String()
		receivedData.Created = created.Format(time.RFC3339)
		message := map[string]interface{}{
			"data": receivedData,
			"type": receivedData.Type,
		}

		for client, id := range activeConnections {
			if id == receivedData.RecipientID {
				err := client.WriteJSON(message)
				if err != nil {
					fmt.Println("Error sending user list to client:", err)
				}

			}

		}
	}
}

func broadcastUserList() {

	var userList []string

	connectionLock.Lock()

	for _, userID := range activeConnections {

		userList = append(userList, userID)

	}
	connectionLock.Unlock()
	// Broadcast the list to all connected clients

	message := map[string]interface{}{
		"data": userList,
		"type": "online-user",
	}
	for client := range activeConnections {

		err := client.WriteJSON(message)
		if err != nil {
			fmt.Println("Error sending user list to client:", err)
		}

	}

}
