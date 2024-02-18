package handlefuncs

// Switch to remarshaling app structs to forward to the frontend.

// Distinguish between errors that need to be returned from, such as
// failure to add item to db, versus errors that just need to be logged,
// such as failure to send message to one of several connections.

// Be consistent about capitalization, e.g. "id" vs "Id" vs "ID".

// Protect the database from concurrent reads and writes. The mattn/go-sqlite3
// documentation says that it's safe for concurrent reads but not for concurrent writes,
// so, as for the activeConnections map, we can use a sync.RWMutex instead of a sync.Mutex,
// as long as we make sure to use .RLock() only when reading from the database and .Lock()
// when writing to it.

// What happens if a user's cookie expires after they've logged in? We need to make sure
// that we're checking for that and handling it appropriately.
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/pkg/db/dbfuncs"
	"sync"
	"time"

	"github.com/google/uuid"
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

func unmarshalBody[T any](signalBody []byte, receivedData T) {
	err := json.Unmarshal(signalBody, receivedData)
	if err != nil {
		log.Println("error unmarshalling body of websocket message:", err)
		log.Println("type of receivedData:", fmt.Sprintf("%T", receivedData))
	}
}

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

type Notification struct {
	Id         string    `json:"Id"`
	ReceiverId string    `json:"RecieverId"`
	SenderId   string    `json:"SenderId"`
	Body       string    `json:"Body"`
	Type       string    `json:"Type"`
	CreatedAt  time.Time `json:"CreatedAt"`
}

type NotificationSeen struct {
	Id string `json:"Id"`
}

// Make methods on this pattern for other app structs as needed.
func (receivedData Notification) parseForDB() *dbfuncs.Notification {
	return &dbfuncs.Notification{
		Body:       receivedData.Body,
		Type:       "requestToFollow",
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
	}
}

// these may not involve database calls but can still be sent through websockets
// this can be resused for sending SignalReceiveds to other users about a user who has
// e.g. registered, logged in, logged out, or changed their status
type BasicUserInfo struct {
	UserId         string `json:"UserId"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Nickname       string `json:"Nickname"`
	PrivacySetting string `json:"PrivacySetting"`
}

type Event struct {
	EventId     string `json:"EventId"`
	GroupId     string `json:"GroupId"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	CreatorId   string `json:"CreatorId"`
}

func (receivedData Event) parseForDB() *dbfuncs.Event {
	return &dbfuncs.Event{
		EventId:     receivedData.EventId,
		GroupId:     receivedData.GroupId,
		Title:       receivedData.Title,
		Description: receivedData.Description,
		CreatorId:   receivedData.CreatorId,
	}
}

func (receivedData Comment) parseForDB() *dbfuncs.Comment {
	return &dbfuncs.Comment{
		Body:      receivedData.Body,
		CreatorId: receivedData.UserID,
		PostId:    receivedData.PostID,
	}
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

camelsBack:
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
				break camelsBack
			}
		}

		err = broker(msgBytes, userID, conn)
		if err.Error() == "logout" {
			break
		}

		var finalStraw error
		if err != nil {
			finalStraw = notifyClientOfError(err, "error processing websocket message", userID)
		}
		if finalStraw != nil {
			log.Println("error sending error message to client:", finalStraw)
			break camelsBack
		}
	}
}

// Check validity of signal and received data as far as possible. But get
// a working version first. We can add more checks later. Checks could
// also be called from the parseForDB methods.
func broker(msgBytes []byte, userID string, conn *websocket.Conn) error {
	var signal SignalReceived
	err := json.Unmarshal(msgBytes, &signal)
	if err != nil {
		log.Println("Error unmarshalling websocket signal:", err)
	}

	switch signal.Type {
	// case "login":
	// No need. This is covered at the start of handleConnection.
	// Sign out and register:
	case "logout":
		logout(userID, conn)
		err = fmt.Errorf("logout")
		return err
	case "notificationSeen":
		var receivedData NotificationSeen
		unmarshalBody(signal.Body, &receivedData)
		err = notificationSeen(receivedData, userID)
	// Chat:
	case "privateMessage":
		var receivedData Message
		unmarshalBody(signal.Body, &receivedData)
		privateMessage(receivedData)
	case "groupMessage":
		var receivedData Message
		unmarshalBody(signal.Body, &receivedData)
		groupMessage(receivedData)
	// General posts, comments, and likes:
	case "post":
		var receivedData Post
		unmarshalBody(signal.Body, &receivedData)
		err = post(receivedData)
	case "comment":
		var receivedData Comment
		unmarshalBody(signal.Body, &receivedData)
		err = comment(receivedData)
	case "like":
		// fill in
	// Group business:
	case "groupCreate":
		//fill in
	case "groupPost":
		var receivedData Post
		unmarshalBody(signal.Body, &receivedData)
		err = groupPost(receivedData)
	case "groupComment":
		//fill in
	case "groupLike":
	// fill in
	// Events:
	case "createEvent":
		var receivedData Event
		unmarshalBody(signal.Body, &receivedData)
		err = createEvent(receivedData)
	case "toggleAttendEvent":
		var receivedData Event
		unmarshalBody(signal.Body, &receivedData)
		// toggleAttendEvent(&receivedData)
		// default:
		// 	//unexpected type
		// 	log.Println("Unexpected websocket message type:", receivedData.Type)
	default:
		// notifications
		var receivedData Notification
		unmarshalBody(signal.Body, &receivedData)
		switch receivedData.Type {
		case "requestToFollow":
			err = requestToFollow(receivedData)
		case "answerRequestToFollow":
			err = answerRequestToFollow(receivedData)
		case "requestToJoinGroup":
			err = requestToJoinGroup(receivedData)
			// case "answerRequestToJoinGroup":
			// 	//fill in
			// case "inviteToJoinGroup":
			// 	//fill in
			// 	// Notify the person you're inviting.
			// 	inviteToJoinGroup(receivedData)
			// case "answerInviteToGroup":
			// 	//fill in
			// 	answerInviteToJoinGroup(receivedData)
		}
	}
	return err
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

// This is a job for web sockets because we want to update the user's
// list of notifications in real time for all their connections, in
// case they have multiple tabs or windows connected.
func notificationSeen(seen NotificationSeen, userID string) error {
	err := dbfuncs.NotificationSeen(seen.Id)
	if err != nil {
		log.Println("error updating notification seen status in database", err)
	}

	connectionLock.RLock()
	for _, c := range activeConnections[userID] {
		err := c.WriteJSON(seen.Id)
		if err != nil {
			log.Println("error sending notification seen status to client", err)
		}
	}

	return err
}

func (receivedData Post) parseForDB() *dbfuncs.Post {
	return &dbfuncs.Post{
		Title:        receivedData.Title,
		Body:         receivedData.Body,
		CreatorId:    receivedData.CreatorId,
		GroupId:      receivedData.GroupId,
		CreatedAt:    receivedData.CreatedAt,
		Image:        receivedData.Image.Data,
		PrivacyLevel: receivedData.PrivacyLevel,
	}
}

func post(receivedData Post) error {
	var err error
	if len(receivedData.Body) > CharacterLimit {
		err = errors.New("413 Payload Too Large")
		log.Println(err)
		return err
	}
	if len(receivedData.Body) == 0 {
		err = errors.New("204 No Content")
		return err
	}

	DBPost := receivedData.parseForDB()
	err = dbfuncs.AddPost(DBPost)
	if err != nil {
		log.Println("error adding post to database", err)
		notifyClientOfError(err, "error adding post to database", receivedData.CreatorId)
		return err
	}

	receivedData.Id, err = uuid.Parse(DBPost.Id)
	if err != nil {
		log.Println("error parsing UUID from database", err)
		return err
	}

	err = send(receivedData)

	return err
}

type PostOrComment interface {
	GetPrivacyLevel() (string, error)
	GetPost() (Post, error)
}

func (receivedData Post) GetPrivacyLevel() (string, error) {
	return receivedData.PrivacyLevel, nil
}

func (receivedData Post) GetPost() (Post, error) {
	return receivedData, nil
}

func (receivedData Comment) GetPrivacyLevel() (string, error) {
	privacyLevel, err := dbfuncs.GetPostPrivacyLevelByCommentId(receivedData.ID)
	return privacyLevel, err
}

func (receivedData Comment) GetPost() (Post, error) {
	dbPost, err := dbfuncs.GetPostByCommentId(receivedData.ID)
	if err != nil {
		log.Println("error getting post from database", err)
		return Post{}, err
	}
	id, err := uuid.Parse(dbPost.Id)
	if err != nil {
		log.Println("error getting post from database", err)
		return Post{}, err
	}
	post := Post{
		Id:           id,
		Title:        dbPost.Title,
		Body:         dbPost.Body,
		CreatorId:    dbPost.CreatorId,
		GroupId:      dbPost.GroupId,
		CreatedAt:    dbPost.CreatedAt,
		Image:        &Image{Data: dbPost.Image},
		PrivacyLevel: dbPost.PrivacyLevel,
	}
	return post, err
}

func send(receivedData PostOrComment) error {
	post, err := receivedData.GetPost()
	if err != nil {
		log.Println("error getting post", err)
		return err
	}
	privacyLevel, err := receivedData.GetPrivacyLevel()
	if err != nil {
		log.Println("error getting privacy level", err)
		return err

	}
	switch privacyLevel {
	case "public":
		connectionLock.RLock()
		for client := range activeConnections {
			for _, c := range activeConnections[client] {
				err := c.WriteJSON(receivedData)
				if err != nil {
					fmt.Println("Error sending user list to client:", err)
					connectionLock.RUnlock()
					return err
				}
			}
		}
		connectionLock.RUnlock()
	case "private":
		followers, err := dbfuncs.GetFollowersByFollowingId(post.CreatorId)
		if err != nil {
			log.Println("error getting followers from database", err)
			notifyClientOfError(err, "error getting followers from database", post.CreatorId)
			return err
		}
		for _, follower := range followers {
			connectionLock.RLock()
			for _, c := range activeConnections[follower] {
				err := c.WriteJSON(receivedData)
				if err != nil {
					fmt.Println("Error sending user list to client:", err)
				}
				connectionLock.RUnlock()
				return err
			}
			connectionLock.RUnlock()
		}
	case "superprivate":
		chosenFollowers, err := dbfuncs.GetPostChosenFollowersByPostId(post.Id.String())
		if err != nil {
			log.Println("error getting chosen followers from database", err)
			notifyClientOfError(err, "error getting chosen followers from database", post.CreatorId)
			return err
		}
		for _, chosen := range chosenFollowers {
			connectionLock.RLock()
			for _, c := range activeConnections[chosen] {
				err := c.WriteJSON(receivedData)
				if err != nil {
					fmt.Println("Error sending user list to client:", err)
				}
				connectionLock.RUnlock()
				return err
			}
			connectionLock.RUnlock()
		}
	}
	connectionLock.RUnlock()

	return err
}

func groupPost(receivedData Post) error {
	var err error
	if len(receivedData.Body) > CharacterLimit {
		err = errors.New("413 Payload Too Large")
		log.Println(err)
		return err
	}
	if len(receivedData.Body) == 0 {
		err = errors.New("204 No Content")
		return err
	}

	DBPost := receivedData.parseForDB()

	err = dbfuncs.AddPost(DBPost)
	if err != nil {
		log.Println("error adding post to database", err)
		notifyClientOfError(err, "error adding post to database", receivedData.CreatorId)
		return err
	}

	recipients, err := dbfuncs.GetGroupMembersByGroupId(receivedData.GroupId)
	if err != nil {
		log.Println("error getting group members from database", err)
		notifyClientOfError(err, "error getting group members from database", receivedData.CreatorId)
		return err
	}

	connectionLock.RLock()
	for _, recipient := range recipients {
		for _, c := range activeConnections[recipient] {
			err := c.WriteJSON(receivedData)
			if err != nil {
				log.Println("error sending group message to recipient", err)
				log.Println("recipient:", recipient)
			}
		}
	}
	connectionLock.RUnlock()

	return err
}

func comment(receivedData Comment) error {
	var err error
	if len(receivedData.Body) > CharacterLimit {
		err = errors.New("413 Payload Too Large")
		log.Println(err)
		return err
	}
	if len(receivedData.Body) == 0 {
		err = errors.New("204 No Content")
		log.Println(err)
		return err
	}

	DBComment := receivedData.parseForDB()
	err = dbfuncs.AddComment(DBComment)
	if err != nil {
		err = errors.New("500 Internal Server Error: error adding comment to database")
		log.Println(err)
		return err
	}

	err = send(receivedData)

	return err
}

// Change this and groupMessage to separate dbfuncs functions. Pass pointer to
// the relevant dbfuncs struct and handle error that's returned.
func privateMessage(receivedData Message) {
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
	if err != nil {
		log.Println("error adding message to database", err)
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
// messages. Update: we've decided have seperate dbfuncs functions
// for private and group messages. TODO: rewrite groupeMessage and
// privateMessage to take account of this.
func groupMessage(receivedData Message) {
	id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
	if err != nil {
		log.Println("error adding message to database", err)
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
				log.Println("error sending group message to recipient", err)
			}
		}
	}
	connectionLock.RUnlock()
}

func notifyClientOfError(err error, message string, id string) error {
	log.Println(err, message)
	data := map[string]interface{}{
		"type": "error",
	}
	connectionLock.RLock()
	for _, c := range activeConnections[id] {
		err = c.WriteJSON(data)
		if err != nil {
			break
		}
	}
	connectionLock.RUnlock()
	return err
}

func requestToFollow(receivedData Notification) error {
	var follow dbfuncs.Follow
	follow.FollowingId = receivedData.ReceiverId
	follow.FollowerId = receivedData.SenderId

	private, err := dbfuncs.IsUserPrivate(receivedData.ReceiverId)
	if err != nil {
		log.Println("error checking if user is public", err)
	}

	if private {
		follow.Status = "pending"
	} else {
		follow.Status = "accepted"
	}

	err = dbfuncs.AddFollow(&follow)
	if err != nil {
		log.Println("error adding follow to database", err)
		notifyClientOfError(err, "error adding follow to database", receivedData.SenderId)
		return err
	}

	if private {
		err = dbfuncs.AddNotification(receivedData.parseForDB())
		if err != nil {
			log.Println("error adding requestToFollow notification to database", err)
			return err
		}

		connectionLock.RLock()
		for _, c := range activeConnections[receivedData.ReceiverId] {
			err := c.WriteJSON(receivedData)
			if err != nil {
				log.Println("error sending group message to recipient", err)
			}
		}
		connectionLock.RUnlock()
	}

	return err
}

func answerRequestToFollow(receivedData Notification) error {
	var err error

	switch receivedData.Body {
	case "yes":
		err = dbfuncs.AcceptFollow(receivedData.SenderId, receivedData.ReceiverId)
		if err != nil {
			log.Println("database error accepting follow", err)
			log.Printf("%s accepted follow request from %s\n",
				receivedData.SenderId, receivedData.ReceiverId)
			return err
		}
	case "no":
		err := dbfuncs.RejectFollow(receivedData.SenderId, receivedData.ReceiverId)
		if err != nil {
			log.Println("error rejecting follow", err)
			log.Printf("%s rejected follow request from %s\n",
				receivedData.SenderId, receivedData.ReceiverId)
		}
		return err
	default:
		log.Println("unexpected body in answerRequestToFollow:", receivedData.Body)
		log.Printf("%s sent unexpected body %s, answering request from %s\n",
			receivedData.SenderId, receivedData.Body, receivedData.ReceiverId)
		return fmt.Errorf("unexpected body in answerRequestToFollow")
	}

	err = dbfuncs.AddNotification(receivedData.parseForDB())
	if err != nil {
		log.Print("error adding notification to database", err)
		log.Printf("%s answered follow request from %s\n",
			receivedData.SenderId, receivedData.ReceiverId)
	}

	connectionLock.RLock()
	for _, c := range activeConnections[receivedData.ReceiverId] {
		err = c.WriteJSON(receivedData)
		if err != nil {
			log.Println("error sending group message to recipient", err)
		}
	}
	connectionLock.RUnlock()

	return err
}

func requestToJoinGroup(receivedData Notification) error {
	notification := dbfuncs.Notification{
		Body:       receivedData.Body,
		Type:       "requestToJoinGroup",
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
		Seen:       false,
	}

	member := dbfuncs.GroupMember{
		GroupId: receivedData.ReceiverId,
		UserId:  receivedData.SenderId,
		Status:  "pending",
	}

	// How about if AddNotification returns the notification ID?
	err := dbfuncs.AddNotification(&notification)
	if err != nil {
		log.Println(err, "error adding notification to database")
		return err
	}
	notificationId := notification.Id
	notification, err = dbfuncs.GetNotificationById(notificationId)
	if err != nil {
		log.Println(err, "error getting notification from database")
		return err
	}

	err = dbfuncs.AddGroupMember(&member)
	if err != nil {
		log.Println(err, "error adding group member to database")
		return err
	}

	creatorId, err := dbfuncs.GetGroupCreatorFromGroupId(receivedData.ReceiverId)
	if err != nil {
		log.Println(err, "error getting group creator from database")
		return err
	}

	connectionLock.RLock()
	for _, c := range activeConnections[creatorId] {
		err := c.WriteJSON(notification)
		if err != nil {
			log.Println(err, "error sending notification to recipient")
		}
	}
	connectionLock.RUnlock()

	return err
}

// // Add a notification to the database and send it to the group creator
// // if they're online. The notification should include the requester and
// // type of notification, the time it was created, and the group ID, in
// // case the recipient has created multiple groups.
// func requestToJoinGroup(receivedData Notification) {
// 	notification := dbfuncs.Notification{
// 		Body:       receivedData.Body,
// 		Type:       "requestToJoinGroup",
// 		ReceiverId: receivedData.RecipientId,
// 		SenderId:   receivedData.SenderId,
// 	}
// 	err := dbfuncs.AddNotification(&notification)
// 	id := notification.Id
// 	if err != nil {
// 		log.Println(err, "error adding notification to database")
// 	}
// 	notification, err = dbfuncs.GetNotificationById(notificationId)
// 	if err != nil {
// 		log.Println(err, "error getting notification from database")
// 	}

// 	connectionLock.RLock()
// 	for _, c := range activeConnections[receivedData.RecipientId] {
// 		err := c.WriteJSON(notification)
// 		if err != nil {
// 			log.Println(err, "error sending notification to recipient")
// 		}
// 	}
// 	connectionLock.RUnlock()
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

// // Add an event to the database and send it to all group members.
// // Add all groups members to GroupEventParticipants with the choice
// // field set to false. Add notification for each group member.
// // In detail: This will now be part of the handlefuncs package.
// // handlefuncs can import dbfuncs, but not the other way around.
// // It will have to convert the Event to a dbfuncs.Event,
// // and then call dbfuncs.AddEvent. It will also have to call
// // dbfuncs.GetGroupMembers to get the list of group members.
// // I've added placeholders for these.
// // It will then have to loop through the group members and call
// // dbfuncs.AddNotification for each one. It will also have to
// // loop through the group members and call dbfuncs.AddGroupEventParticipant
// // for each one. It will also have to loop through the group members
func createEvent(receivedData Event) error {
	event := receivedData.parseForDB()

	id, createdAt, err := dbfuncs.AddEvent(event)
	if err != nil {
		log.Println("error adding event to database", err)
	}

	// make a struct literal with title and description
	body := fmt.Sprintf(`{"Title":%s,"Description":%s}`,
		receivedData.Title, receivedData.Description)

	notification := Notification{
		Id:         id,
		ReceiverId: event.GroupId,
		SenderId:   event.CreatorId,
		Body:       body,
		Type:       "new event",
		CreatedAt:  createdAt,
	}

	connectionLock.RLock()
	recipients := dbfuncs.GetGroupMembers(receivedData.GroupId)
	for _, recipient := range recipients {
		for _, c := range activeConnections[recipient] {
			err := c.WriteJSON(notification)
			if err != nil {
				log.Println("error sending new event notification", err)
			}
		}
	}
	connectionLock.RUnlock()
	return err
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
