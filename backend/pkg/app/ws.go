package app

//general notes/ ideas:
//consider redoing database tables to combine private messages, group messages and various types of notificaitons in one table
//and differentiate them with a type field, possibly with empty fields for fields that are not needed for that type of message.

// Make sure that we're consistent about capitalization, e.g. "id" vs "Id" vs "ID".

// To read from maps, e.g. activeConnections when looping through active connections,
// we can use .RLock() and .RUnlock() instead of .Lock() and .Unlock() to avoid
// blocking other goroutines. I've changes the mutexe types to sync.RWMutex to allow
// this.

// We also need to protect the database from concurrent reads and writes. The mattn/go-sqlite3
// documentation says that it's safe for concurrent reads but not for concurrent writes,
// so, as for the activeConnections map, we can use a sync.RWMutex instead of a sync.Mutex,
// as long as we make sure to use .RLock() only when reading from the database and .RLock()
// when writing to it.

// What happens if a user's cookie expires after they've logged in? We need to make sure
// that we're checking for that and handling it appropriately.

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"server/pkg/db/dbfuncs"
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
	connectionLock    sync.RWMutex
)



// Be sure to allow possibilty of one of a user's connections being closed
// while they still have other connections open. Make this distinct from
// logging out, although logging out will include closing the current
// connection.

// pull some of this stuff out into separate functions to make cleaner
// can consider channel approach instead of mutex
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_token")
	// Add lines to ValidateCookie to RLock the dblock while validating.
	valid := dbfuncs.ValidateCookie(cookie.Value)
	if err != nil || !valid {
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
	}()

	userID, err := dbfuncs.GetUserIdFromCookie(cookie.Value)
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

	// Make sure this works when the user is newly registered. It should
	// broadcast a list including the new user. In realtime, we had the map
	// activeConnections = make(map[*websocket.Conn]string). Here we have
	// activeConnections = make(map[string][]*websocket.Conn). There, when
	// a new user registered, the other users didn't update their lists of
	// users. We need to make sure there is logic in the frontend to handle
	// this.
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

		var signal WsMessage
		err = json.Unmarshal(msgBytes, &signal)
		if err != nil {
			log.Println("Error unmarshalling websocket message:", err)
		}

		switch signal.Type {
		// case "login":
		// No need. This is covered at the start of handleConnection.

		// Sign out and register:
		case "logout":
			handleLogout(userID, conn)
			break eventLoop

		// Chat:
		case "privateMessage":
			var receivedData Message
			unmarshalBody(signal.Body, &receivedData)
			handlePrivateMessage(receivedData)
		case "groupMessage":
			var receivedData Message
			unmarshalBody(signal.Body, &receivedData)
			handleGroupMessage(receivedData)

		// // Cases that require notofications:
		case "requestToFollow":
			// If the request is for a user with a public profile, it should be automatically
			// accepted. Otherwise, notify that user that you want to follow them.
			var receivedData RequestToFollow
			unmarshalBody(signal.Body, &receivedData)
			handleRequestToFollow(receivedData)
		// 	handleRequestToFollow(receivedData)
		// case "answerRequestToFollow":
		// 	answerRequestToFollow(receivedData)
		// case "requestToJoinGroup":
		// 	// Notify the creator.
		// 	requestToJoinGroup(receivedData)
		// case "answerRequestToJoinGroup":
		// 	//fill in
		// 	answerRequestToJoinGroup(receivedData)
		// case "inviteToJoinGroup":
		// 	//fill in
		// 	// Notify the person you're inviting.
		// 	inviteToJoinGroup(receivedData)
		// case "answerInviteToGroup":
		// 	//fill in
		// 	answerInviteToJoinGroup(receivedData)
		case "createEvent":
			var receivedData Event
			unmarshalBody(signal.Body, &receivedData)
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
		case "toggleAttendEvent":
			var receivedData ToggleAttendEvent
			err = json.Unmarshal(signal.Body, &receivedData)
			if err != nil {
				log.Println("Error unmarshalling body of websocket message:", err)
			}
			toggleAttendEvent(&receivedData)

			// default:
			// 	//unexpected type
			// 	log.Println("Unexpected websocket message type:", receivedData.Type)
		}

	}

}

func unmarshalBody[T any](signalBody []byte, receivedData T) {
	err := json.Unmarshal(signalBody, receivedData)
	if err != nil {
		log.Println("Error unmarshalling body of websocket message:", err)
		log.Println("type of receivedData:", fmt.Sprintf("%T", receivedData))
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
	connectionLock.RLock()
	for client := range activeConnections {
		for _, c := range activeConnections[client] {
			err := c.WriteJSON(message)
			if err != nil {
				fmt.Println("Error sending user list to client:", err)
			}
		}
	}
	connectionLock.RUnlock()
}

func closeConnection(conn *websocket.Conn) {
	data := map[string]interface{}{
		"data": "",
		"type": "logout",
	}
	err := conn.WriteJSON(data)
	if err != nil {
		fmt.Println("Error sending logout message to client:", err)
	}
	if err := conn.Close(); err != nil {
		fmt.Println("Error closing websocket:", err)
	}
}

// Tell other connections associated with userID to close themselves
// at the front end. Delete the userID from activeConnections. Broadcast
// the updated user list. The current connection will be closed at the
// when the event loop breaks.
// The frontend also needs to trigger HandleLogOut via http
// as we don't have access to the cookie using websockets.
func handleLogout(userID string, thisConn *websocket.Conn) {
	connectionLock.RLock()
	for _, c := range activeConnections[userID] {
		if c != thisConn {
			closeConnection(c)
		}
	}
	connectionLock.RUnlock()
	connectionLock.Lock()
	delete(activeConnections, userID)
	connectionLock.Unlock()
	broadcastUserList()
}

func handlePrivateMessage(receivedData Message) {
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
	if err != nil {
		log.Println(err, "error adding message to database")
	}
	receivedData.ID = id.String()
	receivedData.Created = created.Format(time.RFC3339)
	message := map[string]interface{}{
		"data": receivedData,
		"type": receivedData.Type,
	}
	connectionLock.RLock()
	for _, c := range activeConnections[receivedData.RecipientID] {
		err := c.WriteJSON(message)
		if err != nil {
			log.Fatal(err)
		}
	}
	connectionLock.RUnlock()
}

// I've adapted dbfuncs.AddMessage to handle both private and group
// messages.
func handleGroupMessage(receivedData Message) {
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
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

	connectionLock.RLock()
	recipients := dbfuncs.GetGroupMembers(receivedData.RecipientID)
	for _, recipient := range recipients {
		for _, c := range activeConnections[recipient] {
			err := c.WriteJSON(message)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	connectionLock.RUnlock()
}

// If the request is for a user with a public profile, add the requester
// to their followers list. Otherwise, add a notification to the database
// and send it to the user with the private profile if they're online.
// (No longer? Also, set the Unseens field of the private user to true.)
// The notification should include the requester and recipient's IDs (and
// maybe usernames), the type of notification, the time it was created, and
// its status (pending, as opposed to accepted or rejected).
func handleRequestToFollow(receivedData RequestToFollow) {
	// // TODO: Write dbfuncs.IsPublic.
	// // Commented out, for now, to avoid red lines.
	// public := dbfuncs.IsPublic(receivedData.RecipientId)
	// if public {
	// 	// Add the requester to the recipient's followers list.
	// } else {
	// 	// Add a notification to the database and send it to the recipient
	// 	// if they're online.
	// }
}

// Decide if we want to notify the requester of the result.
func answerRequestToFollow(receivedData Message) {
}

// Add a notification to the database and send it to the group creator
// if they're online. The notification should include the requester and
// type of notification, the time it was created, and the group ID, in
// case the recipient has created multiple groups.
func requestToJoinGroup(receivedData Message) {
}

// If the answer is yes, add the requester to the group members list
// and broadcast the updated list to all group members. Either way,
// decide if we want to notify the requester of the result. If so,
// and I think we should, add a notification to the database and send
// it to the requester if they're online. The notification should
// include which group and whether the answer was yes or no.
// Also, if YES, add user to GroupEventParticipants with the choice
// field set to false for all events in that group.
func answerRequestToJoinGroup(receivedData Message) {
}

// Add a notification to the database and send it to the person
// being invited if they're online.
func inviteToJoinGroup(receivedData Message) {
}

// If the answer is yes, add the invitee to the group members list
// in the database and broadcast the updated list to all group members.
// Either way, update the status of that notification in the database:
// i.e. change it from pending to accepted or rejected. Decide if we
// want to notify the inviter of the result.
// Also, if YES, add user to GroupEventParticipants with the choice
// field set to false for all events in that group.
func answerInviteToJoinGroup(receivedData Message) {
}

// Add an event to the database and send it to all group members.
// Add all groups members to GroupEventParticipants with the choice
// field set to false. Add notification for each group member.

// In detail: This will now be part of the package.
// can import dbfuncs, but not the other way around.
// It will have to convert the Event to a dbfuncs.Event,
// and then call dbfuncs.AddEvent. It will also have to call
// dbfuncs.GetGroupMembers to get the list of group members.
// It will then have to loop through the group members and call
// dbfuncs.AddNotification for each one. It will also have to
// loop through the group members and call dbfuncs.AddGroupEventParticipant
// for each one. It will also have to loop through the group members

func createEvent(receivedData Event) {
	dbfuncs.AddEvent(receivedData)
	signal := map[string]interface{}{
		"data": receivedData,
		"type": "createEvent",
	}
	connectionLock.RLock()
	recipients := dbfuncs.GetGroupMembers(receivedData.GroupId)
	for _, recipient := range recipients {
		for _, c := range activeConnections[recipient] {
			err := c.WriteJSON(signal)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	connectionLock.RUnlock()
}

// Toggle the user's attendance at the event in the database and
// broadcast the updated list of attendees to all group members
// who are online, or just the change in the user's attendance
// and let the frontend handle the update.
func toggleAttendEvent(receivedData *ToggleAttendEvent) {
	err := dbfuncs.ToggleAttendEvent(receivedData.EventId, receivedData.UserId)
	if err != nil {
		log.Println(err, "error toggling event attendance")
	}

	message := map[string]interface{}{
		"data": receivedData,
		"type": receivedData.Type,
	}

	dbLock.RLock()
	rows, err := db.Query("SELECT * FROM GroupMember WHERE GroupId = ?", receivedData.GroupId)
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
		connectionLock.RLock()
		for _, c := range activeConnections[userId] {
			err := c.WriteJSON(message)
			if err != nil {
				log.Fatal(err)
			}
		}
		connectionLock.RUnlock()
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	dbLock.RUnlock()
}
