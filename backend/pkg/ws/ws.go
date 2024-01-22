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

// Be sure to allow possibilty of one of a user's connections being closed
// while they still have other connections open. Make this distinct from
// logging out, although logging out will include closing the current
// connection.

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

	defer func() {
		closeConnection(conn)
		conn.Close()
	}()

	userID, err := dbfuncs.GetUserIdFromCookie(cookie.Value) // fix spelling of cokie
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

eventLoop:
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
		// case "login":
		// This is covered at the start of handleConnection.

		// Sign out and register:
		case "logout":
			handleLogout(userID)
			break eventLoop
		case "register":
			//fill in
			// Ask why sites make you log in after registering.
			//I.e. Update other users that this user has registered.
			// At present, if the user associated with this connection is an existing
			// user, their arrival is communicated to other users by the call to
			// broadcastUserList before this indefinite for loop. But we also need
			// to communicate when a newly registered users has come online, so that
			// they can be added to the list of users who can be messaged.
		case "updatePrivacySetting":
		//fill in

		// Chat:
		case "privateMessage":
			handlePrivateMessage(receivedData)
		case "groupMessage":
			handleGroupMessage(receivedData)

		// Cases that require notofications:
		case "requestToFollow":
			// If the request is for a user with a public profile, it should be automatically
			// accepted. Otherwise, notify that user that you want to follow them.
			handleRequestToFollow(receivedData)
		case "answerRequestToFollow":
			answerRequestToFollow(receivedData)
		case "requestToJoinGroup":
			// Notify the creator.
			requestToJoinGroup(receivedData)
		case "answerRequestToJoinGroup":
			//fill in
			answerRequestToJoinGroup(receivedData)
		case "inviteToJoinGroup":
			//fill in
			// Notify the person you're inviting.
			inviteToJoinGroup(receivedData)
		case "answerInviteToGroup":
			//fill in
			answerInviteToJoinGroup(receivedData)
		case "createEvent":
			//fill in
			// Notify other group members.
			createEvent(receivedData)

		// General posts, comments, and likes:
		case "post":
			//fill in
		case "comment":
			//fill in
		case "like":
			// fill in

		// Group business:
		case "groupCreate":
			//fill in
		case "groupPost":
			//fill in
		case "groupComment":
			//fill in
		case "groupLike":
		// fill in
		case "attendEvent":
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

// I've left the actual closing of the connection at the backend
// till after this function returns. That way closeConnection can
// be used both to close the other connections associated with
// userID and to close the current connection.
func closeConnection(conn *websocket.Conn) {
	data := map[string]interface{}{
		"data": "",
		"type": "logout",
	}
	err := conn.WriteJSON(data)
	if err != nil {
		fmt.Println("Error sending logout message to client:", err)
	}
	err = conn.Close()
	if err != nil {
		log.Println("Error closing websocket connection:", err)
	}
}

// Tell other connections associated with userID to close themselves
// at the front end. Delete the userID from activeConnections. Broadcast
// the updated user list. The current connection will be closed at the
// when the event loop breaks.
func handleLogout(userID string) {
	for _, c := range activeConnections[userID] {
		closeConnection(c)
		c.Close()
	}
	connectionLock.Lock()
	delete(activeConnections, userID)
	connectionLock.Unlock()
	broadcastUserList()
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

// If the request is for a user with a public profile, add the requester
// to their followers list. Otherwise, add a notification to the database
// and send it to the user with the private profile if they're online.
// Also, set the Unseens field of the private user to true.
// The notification should include the requester and recipient's IDs (and
// maybe usernames), the type of notification, the time it was created, and
// its status (pending, as opposed to accepted or rejected).
func handleRequestToFollow(receivedData handlefuncs.Message) {

}

// If the request is for a user with a public profile, add the requester
// to their followers list. Otherwise, add a notification to the database
// and send it to the user with the private profile if they're online.
// Decide if we want to notify the requester of the result.
func answerRequestToFollow(receivedData handlefuncs.Message) {
}

// Add a notification to the database and send it to the group creator
// if they're online. The notification should include the requester and
// type of notification, the time it was created, and the group ID, in
// case the recipient has created multiple groups.
func requestToJoinGroup(receivedData handlefuncs.Message) {
}

// If the answer is yes, add the requester to the group members list
// and broadcast the updated list to all group members. Either way,
// decide if we want to notify the requester of the result. If so,
// and I think we should, add a notification to the database and send
// it to the requester if they're online. The notification should
// include which group and whether the answer was yes or no.
func answerRequestToJoinGroup(receivedData handlefuncs.Message) {
}

// Add a notification to the database and send it to the person
// being invited if they're online.
func inviteToJoinGroup(receivedData handlefuncs.Message) {
}

// If the answer is yes, add the invitee to the group members list
// in the database and broadcast the updated list to all group members.
// Either way, update the status of that notification in the database:
// i.e. change it from pending to accepted or rejected. Decide if we
// want to notify the inviter of the result.
func answerInviteToJoinGroup(receivedData handlefuncs.Message) {
}

// Add an event to the database and send it to all group members.
func createEvent(receivedData handlefuncs.Message) {
}
