package handlefuncs

//consider redoing database tables to combine private messages, group messages and various types of
// notificaitons in one table
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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/pkg/db/dbfuncs"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:8000"
		},
	}
	// String is userId. The slice of pointers to websocket.Conn is the connections for that user.
	activeConnections = make(map[string][]*websocket.Conn)
	connectionLock    sync.RWMutex
)

// this is what the client sends to the server,
// Body will be unmarshalled based on type into PrivateMessage, GroupMessage, or Notification etc.
// as well as ws messages unrelated to database operations
type SignalReceived struct {
	Type string    `json:"type"`
	Body []byte    `json:"message"`
	Time time.Time `json:"time"`
}

// Consider whether we want to be sending id numbers to the client.
// We can continue for now as if we are, and modify this later if we
// decide not to. Be sure to investigate the security implications
// of sending id numbers to the client. If we decide not to send id
// numbers, we need to consider the actual threat model, and avoid
// naively attempting to fix things that don't actually address the
// real threat.  Matt is sending them and happy that it's not a problem.
// these should match the database fields, check types in database and
// make them match, also move them to db package eventually
// Do they need to match the fields in the database? For example,
// the frontend won't know the id of a message until it's been added
// to the database. At the moment, signals are being sent to the
// frontend in the form of a map from string to interface, with
// whatever key-value pairs are needed for that type of signal.
type PrivateMessage struct {
	Id          string    `json:"Id"`
	SenderId    string    `json:"SenderId"`
	RecipientId string    `json:"RecipientId"`
	Message     string    `json:"Message"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type GroupMessage struct {
	Id        string `json:"Id"`
	SenderId  string `json:"SenderId"`
	GroupId   string `json:"GroupId"`
	Message   string `json:"Message"`
	CreatedAt string `json:"CreatedAt"`
}

// type RequestToFollow struct {
// 	SenderId    string `json:"SenderId"`
// 	RecipientId string `json:"RecipientId"`
// 	Type        string `json:"json"`
// }

type Notification struct {
	RecipientId string    `json:"RecieverId"`
	SenderId    string    `json:"SenderId"`
	Body        string    `json:"Body"`
	Type        string    `json:"Type"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

// these may not involve database calls but can still be sent through websockets
// this can be resused for sending SignalReceiveds to other users about a user who has
// e.g. registered, logged in, logged out, or changed their status
type BasicUserInfo struct {
	UserId         string `json:"UserId"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Nickname       string `json:"Nickname"`
	PrivacySetting string `json:"PrivacySetting"` //maybe?
	//fill in later if more needed
}

type Event struct {
	EventId     string `json:"EventId"`
	GroupId     string `json:"GroupId"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	CreatorId   string `json:"CreatorId"`
}

// Be sure to allow possibilty of one of a user's connections being closed
// while they still have other connections open. Make this distinct from
// logging out, although logging out will include closing the current
// connection.
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_token")
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
	// broadcast a list including the new user, along with their UUID
	// or whatever unique identifier we have. In realtime, we had the map
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
		var signal SignalReceived
		err = json.Unmarshal(msgBytes, &signal)
		if err != nil {
			log.Println("Error unmarshalling websocket message:", err)
		}
		switch signal.Type {
		// case "login":
		// No need. This is covered at the start of handleConnection.
		// Sign out and register:
		case "logout":
			logout(userID, conn)
			break eventLoop
		// Chat:
		case "privateMessage":
			var receivedData Message
			unmarshalBody(signal.Body, &receivedData)
			privateMessage(receivedData)
		case "groupMessage":
			var receivedData Message
			unmarshalBody(signal.Body, &receivedData)
			groupMessage(receivedData)
		// // Cases that require notofications:
		case "requestToFollow":
			var receivedData Notification
			unmarshalBody(signal.Body, &receivedData)
			requestToFollow(receivedData)
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
			var receivedData Event
			unmarshalBody(signal.Body, &receivedData)
			// toggleAttendEvent(&receivedData)
			// default:
			// 	//unexpected type
			// 	log.Println("Unexpected websocket message type:", receivedData.Type)
		}
	}
}

func unmarshalBody[T any](signalBody []byte, receivedData T) {
	err := json.Unmarshal(signalBody, receivedData)
	if err != nil {
		log.Println("error unmarshalling body of websocket message:", err)
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
// The frontend also needs to trigger handlefuncs.HandleLogOut via http
// as we don't have access to the cookie using websockets.
func logout(userID string, thisConn *websocket.Conn) {
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

func privateMessage(receivedData Message) {
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

// I adapted dbfuncs.AddMessage to handle both private and group
// messages. In case, anyone wants to look at it, it should be in
// branch Peter.
func groupMessage(receivedData Message) {
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
	if err != nil {
		log.Println(err, "error adding message to database")
	}
	receivedData.ID = id.String()
	receivedData.Created = created.Format(time.RFC3339)
	message := GroupMessage{
		Id:        "",
		SenderId:  receivedData.SenderID,
		GroupId:   receivedData.RecipientID,
		Message:   receivedData.Message,
		CreatedAt: receivedData.Created,
	}

	connectionLock.RLock()
	recipients := dbfuncs.GetGroupMembers(receivedData.RecipientID)
	for _, recipient := range recipients {
		for _, c := range activeConnections[recipient] {
			err := c.WriteJSON(message)
			if err != nil {
				log.Println(err, "error sending group message to recipient")
			}
		}
	}
	connectionLock.RUnlock()
}

// If the request is for a user with a public profile, add the requester
// to their followers list. Otherwise, add a notification to the database
// and send it to the user with the private profile if they're online.
// The notification should include the requester and recipient's IDs, the
// type of notification, and nickname of requester, the time it was created, and
// its status (pending, as opposed to accepted or rejected). Body of the notification
// should contain the nickname of the requester.
func requestToFollow(receivedData Notification) {
	public, err := dbfuncs.IsPublic(receivedData.RecipientId)
	if err != nil {
		log.Println(err, "error checking if user is public")
	}
	if public {
		var follow dbfuncs.Follow
		follow.FollowingId = receivedData.RecipientId
		follow.FollowerId = receivedData.SenderId
		err := dbfuncs.AddFollower(&follow)
		if err != nil {
			log.Println(err, "error adding follower to database")
		}
	} else {
		notification := dbfuncs.Notification{
			Id:          "",
			Body:        receivedData.Body,
			Type:        "requestToFollow",
			CreatedAt:   time.Now(), // placeholder; let dbfuncs set time while holding the lock
			RecipientId: receivedData.RecipientId,
			SenderId:    receivedData.SenderId,
			Seen:        false,
		}
		err := dbfuncs.AddNotification(&notification)
		if err != nil {
			log.Println(err, "error adding notification to database")
		}
		// Send the notification to the recipient.
		// The notification should include the requester and recipient's IDs, the
		// type of notification, and nickname of requester, the time it was created.
		// Body of the notification should contain the nickname of the requester.

		connectionLock.RLock()
		for _, c := range activeConnections[receivedData.RecipientId] {
			err := c.WriteJSON(notification)
			if err != nil {
				log.Println(err, "error sending notification to recipient")
			}
		}
		connectionLock.RUnlock()
	}
}

// // Decide if we want to notify the requester of the result.
// func answerRequestToFollow(receivedData Message) {
// }

// // Add a notification to the database and send it to the group creator
// // if they're online. The notification should include the requester and
// // type of notification, the time it was created, and the group ID, in
// // case the recipient has created multiple groups.
// func requestToJoinGroup(receivedData Message) {
// }

// // If the answer is yes, add the requester to the group members list
// // and broadcast the updated list to all group members. Either way,
// // decide if we want to notify the requester of the result. If so,
// // and I think we should, add a notification to the database and send
// // it to the requester if they're online. The notification should
// // include which group and whether the answer was yes or no.
// // Also, if YES, add user to GroupEventParticipants with the choice
// // field set to false for all events in that group.
// func answerRequestToJoinGroup(receivedData Message) {
// }

// // Add a notification to the database and send it to the person
// // being invited if they're online.
// func inviteToJoinGroup(receivedData Message) {
// }

// // If the answer is yes, add the invitee to the group members list
// // in the database and broadcast the updated list to all group members.
// // Either way, update the status of that notification in the database:
// // i.e. change it from pending to accepted or rejected. Decide if we
// // want to notify the inviter of the result.
// // Also, if YES, add user to GroupEventParticipants with the choice
// // field set to false for all events in that group.
// func answerInviteToJoinGroup(receivedData Message) {
// }

// Add an event to the database and send it to all group members.
// Add all groups members to GroupEventParticipants with the choice
// field set to false. Add notification for each group member.
// In detail: This will now be part of the handlefuncs package.
// handlefuncs can import dbfuncs, but not the other way around.
// It will have to convert the Event to a dbfuncs.Event,
// and then call dbfuncs.AddEvent. It will also have to call
// dbfuncs.GetGroupMembers to get the list of group members.
// I've added placeholders for these.
// It will then have to loop through the group members and call
// dbfuncs.AddNotification for each one. It will also have to
// loop through the group members and call dbfuncs.AddGroupEventParticipant
// for each one. It will also have to loop through the group members
func createEvent(receivedData Event) {
	dbEvent := dbfuncs.Event{
		EventId:     receivedData.EventId,
		GroupId:     receivedData.GroupId,
		Title:       receivedData.Title,
		Description: receivedData.Description,
		CreatorId:   receivedData.CreatorId,
	}
	err := dbfuncs.AddEvent(&dbEvent)
	if err != nil {
		log.Println(err, "error adding event to database")
	}
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

// // Toggle the user's attendance at the event in the database and
// // broadcast the updated list of attendees to all group members
// // who are online, or just the change in the user's attendance
// // and let the frontend handle the update.

// // Review logic here as it will have changed since our discussions.
// // For example, there will be no direct database access from here.
// func toggleAttendEvent(receivedData *ToggleAttendEvent) {
// 	err := dbfuncs.ToggleAttendEvent(receivedData.EventId, receivedData.UserId)
// 	if err != nil {
// 		log.Println(err, "error toggling event attendance")
// 	}
// 	message := map[string]interface{}{
// 		"data": receivedData,
// 		"type": receivedData.Type,
// 	}
// 	dbLock.RLock()
// 	rows, err := db.Query("SELECT * FROM GroupMember WHERE GroupId = ?", receivedData.GroupId)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var userId string
// 		err = rows.Scan(&userId)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		connectionLock.RLock()
// 		for _, c := range activeConnections[userId] {
// 			err := c.WriteJSON(message)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 		connectionLock.RUnlock()
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	dbLock.RUnlock()
// }
