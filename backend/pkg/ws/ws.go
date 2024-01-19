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

	// String is userID. The slice of pointers to websocket.Conn is the connections for that user.
	activeConnections = make(map[string][]*websocket.Conn)
	connectionLock    sync.Mutex
	dbLock            sync.Mutex
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
func HandleConnection(w http.ResponseWriter, r *http.Request) {
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
	if _, ok := activeConnections[userID]; !ok {
		activeConnections[userID] = []*websocket.Conn{conn}
	} else {
		activeConnections[userID] = append(activeConnections[userID], conn)
	}
	connectionLock.Unlock()

	broadcastUserList()

	for {
		_, msgBytes, err := conn.ReadMessage()
		//possibly don't want to immediately delete the connection if there is an error
		if err != nil {
			myUpdatedConnections := []*websocket.Conn{}
			for _, c := range activeConnections[userID] {
				if c != conn {
					myUpdatedConnections = append(myUpdatedConnections, conn)
				}
				connectionLock.Lock()
				activeConnections[userID] = myUpdatedConnections
				if len(myUpdatedConnections) == 0 {
					delete(activeConnections, userID)
					log.Println("User", userID, "disconnected, unable to read from websocket, error:", err)
				}
				connectionLock.Unlock()
				break
			}
		}

		var receivedData handlefuncs.Message
		err = json.Unmarshal(msgBytes, &receivedData)
		if err != nil {
			log.Println("Error unmarshalling websocket message:", err)
		}

		switch receivedData.Type {
		case "privateMessage":
			handlePrivateMessage(receivedData)
		case "groupMessage":
			handleGroupMessage(receivedData)
		case "requestToFollow":
			//fill in
			// If the request is for a user with a public profile, it should be automatically
			// accepted. Otherwise, notify that user that you want to follow them.
		case "answerRequestToFollow":
			//fill in
		case "requestToJoinGroup":
			//fill in
			// Notify the creator.
		case "answerRequestToJoinGroup":
			//fill in
		case "inviteToGroup":
			//fill in
			// Notify the person you're inviting.
		case "answerInviteToGroup":
			//fill in
		case "createEvent":
			//fill in
			// Notify other group members.
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
			//I.e. Update other users that this user has registered.
			// At present, if the user associated with this connection is an existing
			// user, their arrival is communicated to other users by the call to
			// broadcastUserList before this indefinite for loop. But we also need
			// to communicate when a newly registered users has come online, so that
			// they can be added to the list of users who can be messaged.
		case "updatePrivacySetting":
			//fill in
		default:
			//unexpected type
			log.Println("Unexpected websocket message type:", receivedData.Type)
		}

	}

}

func broadcastUserList() {
	var userList []string

	connectionLock.Lock()
	for userID := range activeConnections {
		userList = append(userList, userID)
	}
	connectionLock.Unlock()

	message := map[string]interface{}{
		"data": userList,
		"type": "online-user",
	}

	for client := range activeConnections {
		for _, c := range activeConnections[client] {
			err := c.WriteJSON(message)
			if err != nil {
				fmt.Println("Error sending user list to client:", err)
			}
		}
	}

}

func handlePrivateMessage(receivedData handlefuncs.Message) {
	dbLock.Lock()
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
	dbLock.Unlock()
	if err != nil {
		log.Println(err, "error adding message to database")
	}
	receivedData.ID = id.String()
	receivedData.Created = created.Format(time.RFC3339)
	message := map[string]interface{}{
		"data": receivedData,
		"type": receivedData.Type,
	}
	for _, c := range activeConnections[receivedData.RecipientID] {
		err := c.WriteJSON(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleGroupMessage(receivedData handlefuncs.Message) {
	dbLock.Lock()
	// Change this to AddGroupMessage or adapt AddMessage to handle both.
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
	dbLock.Unlock()
	if err != nil {
		log.Println(err, "error adding message to data base")
	}
	receivedData.ID = id.String()
	receivedData.Created = created.Format(time.RFC3339)
	message := map[string]interface{}{
		"data":    receivedData,
		"type":    receivedData.Type,
		"groupId": receivedData.RecipientID,
	}

	rows, err := db.Query("SELECT * FROM GroupMembers WHERE GroupId = ?", receivedData.RecipientID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range activeConnections[userId] {
			err := c.WriteJSON(message)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
