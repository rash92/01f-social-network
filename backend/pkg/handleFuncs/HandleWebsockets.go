package handlefuncs

// Separete out Notifications type. I was thinking of taking a short-cut
// and doing all the notifications as one type, but now I think it's
// better to separate them.

//

// Some errors sending signals to clients might just be because they've
// gone offline, and that's fine.

// You can only have one cookie per browser, but you can have multiple
// connections per cookie. This is because the cookie is stored in the
// browser, and the browser is the client. The connection is stored on
// the server.

// Include LastMessageTime in database for users so that we can
// order them in the chat. And include logic for that here.

// If client receives a signal of type BasicUserInfo, that means they
// have a (potential) new follower.

// cookies, connections, and websockets: multiple cookies per id?
// In the original forum, we had one cookie per user. When someone
// logged in, we deleted any existing cookie and gave them a new one.
// Here we'd like to allow users to be logged in on multiple devices
// and browsers, as well as just in multiple tabs. To do this, we could
// just not delete the cookie, or we could delete it only when the user
// is not already logged in. This means we'd need to allow for the
// possibility of multiple cookies per user. That's fine. The database
// already allows for this.

// *

// Switching to remarshaling app structs to forward to the frontend.
// I'm up to requestToJoinGroup with this.

// Changing uuid.UUID to string everywhere as I go along.

// Note that notifyClientOfError only notifies the client that there was
// an error. It doesn't tell the client what the error was. From the
// POV of the client, it will always be an internal server error. Nor
// should it notify a user if someone else failed in some action that
// would have affected them only if it had succeeded.

// Distinguish between errors that need to be returned from, such as
// failure to add item to db, versus errors that just need to be logged,
// such as failure to send message to one of several connections.

// Be consistent about capitalization, e.g. "id" vs "Id" vs "ID". I'm
// changing app structs to match database column names as I go along,
// thus "Id" rather than "ID".

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
	Type string `json:"type"`
	Body []byte `json:"message"`
}

func unmarshalBody[T any](signalBody []byte, receivedData T) {
	err := json.Unmarshal(signalBody, receivedData)
	if err != nil {
		log.Println("error unmarshalling body of websocket message:", err)
		log.Println("type of receivedData:", fmt.Sprintf("%T", receivedData))
	}
}

type Follow struct {
	FollowerId  string `json:"FollowrId"`
	FollowingId string `json:"FollowingId"`
	Status      string `json:"Status"`
}

type PrivateMessage struct {
	Id         string    `json:"Id"`
	SenderId   string    `json:"SenderId"`
	ReceiverId string    `json:"ReceiverId"`
	Message    string    `json:"Message"`
	CreatedAt  time.Time `json:"CreatedAt"`
}

type GroupMessage struct {
	Id        string    `json:"Id"`
	SenderId  string    `json:"SenderId"`
	GroupId   string    `json:"GroupId"`
	Message   string    `json:"Message"`
	CreatedAt time.Time `json:"CreatedAt"`
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

type Group struct {
	Id          string    `json:"Id"`
	CreatorId   string    `json:"CreatorId"`
	Title       string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type InviteToJoinGroup struct {
	GroupId  string `json:"GroupId"`
	UserId   string `json:"UserId"`
	SenderId string `json:"Sender"`
}

type Notification struct {
	Id         string    `json:"Id"`
	ReceiverId string    `json:"RecieverId"`
	SenderId   string    `json:"SenderId"`
	Body       string    `json:"Body"`
	Type       string    `json:"Type"`
	CreatedAt  time.Time `json:"CreatedAt"`
	Seen       bool      `json:"Seen"`
}

type NotificationSeen struct {
	Id string `json:"Id"`
}

// Need to keep Notification type general, not just requestToFollow!
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

type RequestToFollow struct {
	User   BasicUserInfo `json:"User"`
	Status string        `json:"Status"`
	Type   string        `json:"Type"`
}

type GroupMember struct {
	GroupId string `json:"GroupId"`
	UserId  string `json:"UserId"`
	Status  string `json:"Status"`
}

type AnswerRequestToJoinGroup struct {
	SenderId   string `json:"SenderId"`
	ReceiverId string `json:"ReceiverId"`
	GroupId    string `json:"GroupId"`
	Accept     bool   `json:"Accept"`
}

type GroupEvent struct {
	Id          string    `json:"EventId"`
	GroupId     string    `json:"GroupId"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	CreatorId   string    `json:"CreatorId"`
	Time        time.Time `json:"Time"`
}

func (receivedData GroupEvent) parseForDB() *dbfuncs.GroupEvent {
	return &dbfuncs.GroupEvent{
		Id:          receivedData.Id,
		GroupId:     receivedData.GroupId,
		Title:       receivedData.Title,
		Description: receivedData.Description,
		CreatorId:   receivedData.CreatorId,
	}
}

type GroupEventParticipant struct {
	EventId string `json:"EventId"`
	UserId  string `json:"UserId"`
	GroupId string `json:"GroupId"`
	Choice  string `json:"Choice"`
}

func (receivedData Comment) parseForDB() *dbfuncs.Comment {
	return &dbfuncs.Comment{
		Body:      receivedData.Body,
		CreatorId: receivedData.UserId,
		PostId:    receivedData.PostId,
	}
}

// When a new user registers, send their basic info to all other users? No,
// this was my original thought, but that would only be necessary if we
// displayed a menu of all users to all users, as we did in the realtime
// chat app. We don't do that here. We only display a list of users who
// are followers or following.

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

	userID, err := dbfuncs.GetUserIdFromCookie(cookie.Value)
	if err != nil {
		log.Println("Error retrieving userID from database:", err)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading to WebSocket:", err)
		return
	}
	defer func() {
		conn.Close()
	}()

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
		if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Printf("error: %v", err)
			myUpdatedConnections := []*websocket.Conn{}
			connectionLock.Lock()
			_, ok := activeConnections[userID]
			if ok {
				for _, c := range activeConnections[userID] {
					if c != conn {
						myUpdatedConnections = append(myUpdatedConnections, conn)
					}
				}
				if len(myUpdatedConnections) == 0 {
					delete(activeConnections, userID)
				}
				activeConnections[userID] = myUpdatedConnections
			}
			connectionLock.Unlock()
			break
		}

		err = broker(msgBytes, userID, conn, w, r)
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
func broker(msgBytes []byte, userID string, conn *websocket.Conn, w http.ResponseWriter, r *http.Request) error {
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
		HandleLogout(w, r)
		return err
	case "requestToFollow":
		var receivedData Follow
		unmarshalBody(signal.Body, &receivedData)
		err = requestToFollow(receivedData)
	case "answerRequestToFollow":
		var receivedData Notification
		err = answerRequestToFollow(receivedData)
	case "notificationSeen":
		var receivedData NotificationSeen
		unmarshalBody(signal.Body, &receivedData)
		err = notificationSeen(receivedData, userID)

	// Chat:
	case "privateMessage":
		var receivedData PrivateMessage
		unmarshalBody(signal.Body, &receivedData)
		err = privateMessage(receivedData)
	case "groupMessage":
		var receivedData GroupMessage
		unmarshalBody(signal.Body, &receivedData)
		err = groupMessage(receivedData)

	// General posts, comments, and likes:
	case "post":
		var receivedData Post
		unmarshalBody(signal.Body, &receivedData)
		err = postOrComment(receivedData)
	case "comment":
		var receivedData Comment
		unmarshalBody(signal.Body, &receivedData)
		err = postOrComment(receivedData)
	case "like":
		var receievedData Like
		unmarshalBody(signal.Body, &receievedData)
		err = like(receievedData)

		// Groups:
	case "createGroup":
		var receievedData Group
		unmarshalBody(signal.Body, &receievedData)
		err = createGroup(receievedData)
	case "answerRequestToJoinGroup":
		var receivedData AnswerRequestToJoinGroup
		unmarshalBody(signal.Body, &receivedData)
		err = answerRequestToJoinGroup(receivedData)
	case "groupPost":
		var receivedData Post
		unmarshalBody(signal.Body, &receivedData)
		err = postOrComment(receivedData)
	case "groupComment":
		var receivedData Comment
		unmarshalBody(signal.Body, &receivedData)
		err = postOrComment(receivedData)

	// Events:
	case "requestToJoinGroup":
		var receivedData GroupMember
		unmarshalBody(signal.Body, &receivedData)
		err = requestToJoinGroup(receivedData)
	case "inviteToJoinGroup":
		var receivedData InviteToJoinGroup
		unmarshalBody(signal.Body, &receivedData)
		err = inviteToJoinGroup(receivedData)
	case "createEvent":
		var receivedData GroupEvent
		unmarshalBody(signal.Body, &receivedData)
		err = createEvent(receivedData)
	case "toggleAttendEvent":
		var receivedData GroupEventParticipant
		unmarshalBody(signal.Body, &receivedData)
		toggleAttendEvent(receivedData)
	default:
		err = errors.New("unexpected websocket message type")
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
}

// Tell other connections associated with userID to close themselves
// at the front end. Delete the userID from activeConnections. Broadcast
// the updated user list. The current connection will be closed at the
// when the event loop breaks.
// The frontend also needs to trigger handlefuncs.HandleLogOut via http
// as we don't have access to the cookie using websockets. This needs
// to be done by each open connection, as well as by the one that
// initiated the logout.
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

func validateContent(content string) error {
	if len(content) > CharacterLimit {
		err := errors.New("413 Payload Too Large")
		log.Println(err)
		return err
	}
	if len(content) == 0 {
		err := errors.New("204 No Content")
		log.Println(err)
		return err
	}
	return nil
}

func postOrComment(receivedData PostOrComment) error {
	err := validateContent(receivedData.GetBody())
	if err != nil {
		return err
	}

	err = receivedData.SetId()
	if err != nil {
		return err
	}

	err = send(receivedData)
	return err
}

type PostOrComment interface {
	GetPrivacyLevel() (string, error)
	GetPost() (Post, error)
	GetBody() string
	GetId() string
	AddToDB() (string, error)
	GetCreatorId() string
	SetId() error
}

// The following methods are just to satisfy the PostOrComment interface,
// which I defined so that I could pass Posts and Comments to one
// postOrComment function. They may seem redundant, but I think they
// do save some repetition, as they make it possible to have a single
// postOrComment function that can handle posts, group posts, and comments.
func (receivedData Post) GetBody() string {
	return receivedData.Body
}

func (receivedData Comment) GetBody() string {
	return receivedData.Body
}

func (receivedData Post) GetCreatorId() string {
	return receivedData.CreatorId
}

func (receivedData Comment) GetCreatorId() string {
	return receivedData.UserId
}

func (receivedData Post) GetId() string {
	return receivedData.Id
}

func (receivedData Comment) GetId() string {
	return receivedData.Id
}

func (receivedData Post) SetId() error {
	var err error
	receivedData.Id, err = receivedData.AddToDB()
	if err != nil {
		log.Println("error adding content to database", err)
		notifyClientOfError(err, "error adding content to database", receivedData.GetCreatorId())
	}
	return err
}

func (receivedData Comment) SetId() error {
	var err error
	receivedData.Id, err = receivedData.AddToDB()
	if err != nil {
		log.Println("error adding content to database", err)
		notifyClientOfError(err, "error adding content to database", receivedData.GetCreatorId())
	}
	return err
}

func (receivedData Post) GetPrivacyLevel() (string, error) {
	return receivedData.PrivacyLevel, nil
}

func (receivedData Post) GetPost() (Post, error) {
	return receivedData, nil
}

func (receivedData Comment) GetPrivacyLevel() (string, error) {
	privacyLevel, err := dbfuncs.GetPostPrivacyLevelByCommentId(receivedData.Id)
	return privacyLevel, err
}

func (receivedData Comment) GetPost() (Post, error) {
	postId, err := dbfuncs.GetPostIdByCommentId(receivedData.Id)
	if err != nil {
		log.Println("error getting post id from database", err)
		return Post{}, err
	}

	dbPost, err := dbfuncs.GetPostById(postId)
	if err != nil {
		log.Println("error getting post from database", err)
		return Post{}, err
	}

	post := Post{
		Id:           dbPost.Id,
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

	var privacyLevel string
	if post.GroupId != "" {
		privacyLevel = "group"
	} else {
		privacyLevel, err = receivedData.GetPrivacyLevel()
		if err != nil {
			log.Println("error getting privacy level", err)
			return err
		}
	}

	switch privacyLevel {
	case "group":
		connectionLock.RLock()
		recipients := dbfuncs.GetGroupMembers(post.GroupId)
		for _, recipient := range recipients {
			for _, c := range activeConnections[recipient] {
				err := c.WriteJSON(receivedData)
				if err != nil {
					err = errors.New("error sending group message to recipient")
					log.Println(err)
				}
			}
		}
		connectionLock.RUnlock()
		return err
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
		followers, err := dbfuncs.GetAcceptedFollowerIdsByFollowingId(post.CreatorId)
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
		chosenFollowers, err := dbfuncs.GetPostChosenFollowerIdsByPostId(post.Id)
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

func (receivedData Comment) AddToDB() (string, error) {
	dbStruct := receivedData.parseForDB()
	err := dbfuncs.AddComment(dbStruct)
	if err != nil {
		return "", err
	}
	return dbStruct.Id, err
}

func (receivedData Post) AddToDB() (string, error) {
	dbStruct := receivedData.parseForDB()
	err := dbfuncs.AddPost(dbStruct)
	if err != nil {
		return "", err
	}
	return dbStruct.Id, err
}

func privateMessage(receivedData PrivateMessage) error {
	dbPM := dbfuncs.PrivateMessage{
		SenderId:   receivedData.SenderId,
		ReceiverId: receivedData.ReceiverId,
		Message:    receivedData.Message,
	}

	err := dbfuncs.AddPrivateMessage(&dbPM)
	if err != nil {
		log.Println("error adding message to database", err)
		notifyClientOfError(err, "error adding message to database", receivedData.SenderId)
		return err
	}

	receivedData.Id = dbPM.Id
	receivedData.CreatedAt = dbPM.CreatedAt

	connectionLock.RLock()
	defer connectionLock.RUnlock()

	for _, c := range activeConnections[receivedData.ReceiverId] {
		err := c.WriteJSON(receivedData)
		if err != nil {
			log.Println("error sending private message to recipient", err)
		}
	}

	return err
}

func groupMessage(receivedData GroupMessage) error {
	dbGM := dbfuncs.GroupMessage{
		SenderId: receivedData.SenderId,
		GroupId:  receivedData.GroupId,
		Message:  receivedData.Message,
	}

	err := dbfuncs.AddGroupMessage(&dbGM)
	if err != nil {
		log.Println("error adding message to database", err)
		notifyClientOfError(err, "error adding message to database", receivedData.SenderId)
	}

	receivedData.Id = dbGM.Id
	receivedData.CreatedAt = dbGM.CreatedAt

	connectionLock.RLock()
	defer connectionLock.RUnlock()

	recipients, err := dbfuncs.GetGroupMemberIdsByGroupId(receivedData.GroupId)
	if err != nil {
		log.Println("error getting group members from database", err)
		notifyClientOfError(err, "error getting group members from database", receivedData.SenderId)
		return err
	}

	for _, recipient := range recipients {
		for _, c := range activeConnections[recipient] {
			err := c.WriteJSON(receivedData)
			if err != nil {
				log.Println("error sending group message to recipient", err)
			}
		}
	}

	return err
}

// Only notify a user of an error that occurred while processing an
// action they attempted. No need to notify someone if someone else
// failed to follow them, for example. I'm thinking, also, only notify
// user that a message couldn't be added to the db, since that affects
// them directly. If a message couldn't be sent to one of their connections,
// we can just log that and deal with it ourselves.
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

func requestToFollow(receivedData Follow) error {
	var follow dbfuncs.Follow
	follow.FollowingId = receivedData.FollowingId
	follow.FollowerId = receivedData.FollowerId

	private, err := dbfuncs.IsUserPrivate(receivedData.FollowingId)
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
		notifyClientOfError(err, "error adding follow to database", receivedData.FollowerId)
		return err
	}

	follower, err := dbfuncs.GetUserById(receivedData.FollowerId)
	if err != nil {
		log.Println("error getting nickname from database", err)
		notifyClientOfError(err, "error getting nickname from database", receivedData.FollowerId)
		return err
	}

	notification := Notification{
		ReceiverId: receivedData.FollowingId,
		SenderId:   receivedData.FollowerId,
		Body:       fmt.Sprintf("%s has requested to follow you", follower.Nickname),
		Type:       "requestToFollow",
	}

	if private {
		err = dbfuncs.AddNotification(notification.parseForDB())
		if err != nil {
			log.Println("error adding requestToFollow notification to database", err)
			return err
		}
	}

	connectionLock.RLock()
	for _, c := range activeConnections[follow.FollowingId] {
		err = c.WriteJSON(notification)
		if err != nil {
			log.Println("error sending (potential) new follower info to recipient", err)
		}
	}
	connectionLock.RUnlock()

	return err
}

// When received, client should request profile if they're on profile page.
func answerRequestToFollow(receivedData Notification) error {
	var err error

	switch receivedData.Body {
	case "yes":
		err = dbfuncs.AcceptFollow(receivedData.SenderId, receivedData.ReceiverId)
		if err != nil {
			log.Println("database error accepting follow", err)
			notifyClientOfError(err, "database error accepting follow", receivedData.SenderId)
			return err
		}
	case "no":
		err := dbfuncs.DeleteFollow(receivedData.SenderId, receivedData.ReceiverId)
		if err != nil {
			log.Println("error rejecting follow", err)
			notifyClientOfError(err, "error rejecting follow", receivedData.SenderId)
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

	receivedData.Seen = false

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

func createGroup(receivedData Group) error {
	group := dbfuncs.Group{
		CreatorId:   receivedData.CreatorId,
		Title:       receivedData.Title,
		Description: receivedData.Description,
	}

	err := dbfuncs.AddGroup(&group)
	if err != nil {
		log.Println("error adding group to database", err)
		notifyClientOfError(err, "error adding group to database", receivedData.CreatorId)
		return err
	}

	receivedData.Id = group.Id
	receivedData.CreatedAt = group.CreatedAt

	member := dbfuncs.GroupMember{
		GroupId: group.Id,
		UserId:  receivedData.CreatorId,
		Status:  "accepted",
	}

	err = dbfuncs.AddGroupMember(&member)
	if err != nil {
		log.Println("error adding group member to database", err)
		notifyClientOfError(err, "error adding group member to database", receivedData.CreatorId)
		return err
	}

	connectionLock.RLock()
	for user := range activeConnections {
		for _, c := range activeConnections[user] {
			err = c.WriteJSON(group)
			if err != nil {
				log.Println("error sending new group to client", err)
			}
		}
	}
	connectionLock.RUnlock()

	return err
}

func requestToJoinGroup(receivedData GroupMember) error {
	groupCreator, err := dbfuncs.GetGroupCreatorIdByGroupId(receivedData.GroupId)
	if err != nil {
		log.Println("error finding group creator")
		notifyClientOfError(err, "error finding group creator", receivedData.UserId)
		return err
	}

	member := dbfuncs.GroupMember{
		GroupId: receivedData.GroupId,
		UserId:  receivedData.UserId,
		Status:  "pending",
	}

	err = dbfuncs.AddGroupMember(&member)
	if err != nil {
		log.Println(err, "error adding new group member to database")
		notifyClientOfError(err, "error adding new group member to database", receivedData.UserId)
		return err
	}

	dbNotification := dbfuncs.Notification{
		Type:       "requestToJoinGroup",
		ReceiverId: groupCreator,
		SenderId:   receivedData.UserId,
	}

	err = dbfuncs.AddNotification(&dbNotification)
	if err != nil {
		log.Println(err, "error adding notification to database")
		notifyClientOfError(err, "error adding notification to database", receivedData.UserId)
		return err
	}

	notification := Notification{
		Id:         dbNotification.Id,
		Type:       "requestToJoinGroup",
		ReceiverId: groupCreator,
		SenderId:   receivedData.UserId,
	}

	connectionLock.RLock()
	for _, c := range activeConnections[groupCreator] {
		err := c.WriteJSON(notification)
		if err != nil {
			log.Println(err, "error sending notification to recipient")
		}
	}
	connectionLock.RUnlock()

	return err
}

func answerRequestToJoinGroup(receivedData AnswerRequestToJoinGroup) error {
	var err error

	member := dbfuncs.GroupMember{
		GroupId: receivedData.GroupId,
		UserId:  receivedData.ReceiverId,
		Status:  "accepted",
	}

	if receivedData.Accept {
		err = dbfuncs.UpdateGroupMember(&member)
		if err != nil {
			log.Println(err, "error updating group member status in database")
			notifyClientOfError(err, "error updating group member status in database", receivedData.SenderId)
			return err
		}
	} else {
		err = dbfuncs.DeleteGroupMember(&member)
		if err != nil {
			log.Println(err, "error deleting group member from database")
			notifyClientOfError(err, "error deleting group member from database", receivedData.SenderId)
			return err
		}
	}

	dbNotification := dbfuncs.Notification{
		Type:       "answerRequestToJoinGroup",
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
	}

	err = dbfuncs.AddNotification(&dbNotification)
	if err != nil {
		log.Println(err, "error adding notification to database")
	}

	notification := Notification{
		Id:         dbNotification.Id,
		Type:       "answerRequestToJoinGroup",
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
	}

	if receivedData.Accept {
		notification.Body = "yes"
	} else {
		notification.Body = "no"
	}

	connectionLock.RLock()
	for _, c := range activeConnections[receivedData.ReceiverId] {
		err := c.WriteJSON(notification)
		if err != nil {
			log.Println(err, "error sending notification to recipient")
		}
	}
	connectionLock.RUnlock()

	return err
}

// Re errors, what is our threat model here? Suppose the pending user
// is added to the GroupMembers table, but the notification fails to
// be added to the Notifications table. Should we delete the user from
// the GroupMembers table and notify the sender that it failer? We can't
// addd the pending member but not notify them. But then what if the
// notification of the error fails? The server could have a collection
// of failed notifications to keep trying to add to the database or
// to send or whatever. I suppose such measures might catch some errors,
// if not all. But we also need to finish our project ... There is
// the danger that the server could be overwhelmed by failed notifications,
// and the danger that the database could be overwhelmed by attemots to
// add the same entry again and again, if it's rightly OR WRONGLY
// identified as a failed attempt! This could get expensive.
func inviteToJoinGroup(receivedData InviteToJoinGroup) error {
	member := dbfuncs.GroupMember{
		GroupId: receivedData.GroupId,
		UserId:  receivedData.UserId,
		Status:  "pending",
	}

	err := dbfuncs.AddGroupMember(&member)
	if err != nil {
		log.Println(err, "error adding new group member to database")
		notifyClientOfError(err, "error adding new group member to database", receivedData.SenderId)
		return err
	}

	dbNotification := dbfuncs.Notification{
		Type:       "inviteToJoinGroup",
		ReceiverId: receivedData.UserId,
		SenderId:   receivedData.SenderId,
	}

	sender, err := dbfuncs.GetUserById(receivedData.SenderId)
	if err != nil {
		log.Println(err, "error getting sender info from database")
		return err
	}

	dbNotification.Body = fmt.Sprintf("%s invited you to join group %s", sender, receivedData.GroupId)

	err = dbfuncs.AddNotification(&dbNotification)
	if err != nil {
		log.Println(err, "error adding notification to database")
		notifyClientOfError(err, "error adding notification to database", sender.Nickname)
		return err
	}

	notification := Notification{
		Id:         dbNotification.Id,
		Type:       "inviteToJoinGroup",
		ReceiverId: receivedData.UserId,
		SenderId:   receivedData.SenderId,
	}

	connectionLock.RLock()
	for _, c := range activeConnections[receivedData.UserId] {
		err := c.WriteJSON(notification)
		if err != nil {
			log.Println(err, "error sending notification to recipient")
		}
	}
	connectionLock.RUnlock()

	return err
}

// // If the answer is yes, add the invitee to the group members list
// // in the database and broadcast the updated list to all group members.
// // Either way, update the status of that notification in the database:
// // i.e. change it from pending to accepted or rejected. Decide if we
// // want to notify the inviter of the result.
// // Also, if YES, add user to GroupEventParticipants with the choice
// // field set to false for all events in that group.
// func answerInviteToJoinGroup(receivedData Message) {
// }

// I'm assuming we'll have a notification for each group member.
// The alternative would be to have one for the group. This way
// means more notifications, but less logic.
func createEvent(receivedData GroupEvent) error {
	event := receivedData.parseForDB()

	err := dbfuncs.AddGroupEvent(event)
	if err != nil {
		log.Println("error adding event to database", err)
		notifyClientOfError(err, "error adding event to database", receivedData.CreatorId)
		return err
	}

	body := fmt.Sprintf(`{"Title":%s,"Description":%s}`,
		receivedData.Title, receivedData.Description)

	members, err := dbfuncs.GetGroupMembersByGroupId(receivedData.GroupId)
	if err != nil {
		log.Println("error getting group members from database", err)
	}

	for _, member := range members {
		notification := Notification{
			ReceiverId: member,
			SenderId:   receivedData.CreatorId,
			Body:       body,
			Type:       "event",
		}
		connectionLock.RLock()
		for _, c := range activeConnections[member] {
			err := c.WriteJSON(notification)
			if err != nil {
				log.Println("error sending new event notification", err)
			}
		}
		connectionLock.RUnlock()
	}

	return err
}

func (receivedData GroupEventParticipant) parseForDB() dbfuncs.GroupEventParticipant {
	return dbfuncs.GroupEventParticipant{
		EventId: receivedData.EventId,
		UserId:  receivedData.UserId,
		GroupId: receivedData.GroupId,
		Choice:  receivedData.Choice,
	}
}

func toggleAttendEvent(receivedData GroupEventParticipant) error {
	participant := receivedData.parseForDB()

	err := dbfuncs.ToggleAttendEvent(&participant)
	if err != nil {
		log.Println("error toggling event attendance in database", err)
		notifyClientOfError(err, "error toggling event attendance in database", receivedData.UserId)
		return err
	}

	members, err := dbfuncs.GetGroupMemberIdsByGroupId(receivedData.GroupId)
	if err != nil {
		log.Println("error getting group members from database", err)
		return err
	}

	for _, member := range members {
		connectionLock.RLock()
		for _, c := range activeConnections[member] {
			err := c.WriteJSON(receivedData)
			if err != nil {
				log.Println("error sending event attendance update to client", err)
			}
		}
		connectionLock.RUnlock()
	}

	return err
}

// LikeDislikePost(UserId, PostId, likeOrDislike string)

type Like struct {
	UserId        string `json:"UserId"`
	PostId        string `json:"PostId"`
	LikeOrDislike string `json:"LikeOrDislike"`
}

func like(receivedData Like) error {
	err := dbfuncs.LikeDislikePost(receivedData.UserId, receivedData.PostId, receivedData.LikeOrDislike)

	// Get privacy level of post or comment by analogy, then send to
	// appropriate recipients by analogy with send().

	return err
}
