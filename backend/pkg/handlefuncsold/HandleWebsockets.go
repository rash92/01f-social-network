package handlefuncs

// Include LastMessageTime in database for users so that we can
// order them in the chat. And include logic for that here.

// If client receives a signal of type BasicUserInfo, that means they
// have a (potential) new follower.

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
	"backend/pkg/db/dbfuncs"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
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
			return origin == "http://localhost:3000"
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
	Type string          `json:"type"`
	Body json.RawMessage `json:"message"`
}

func unmarshalBody[T any](signalBody []byte, receivedData T) {
	err := json.Unmarshal(signalBody, receivedData)
	if err != nil {
		log.Println("error unmarshalling body of websocket message:", err)
		log.Println("type of receivedData:", fmt.Sprintf("%T", receivedData))
	}
}

// type PrivateMessage struct {
// 	Id          string    `json:"Id"`
// 	SenderId    string    `json:"SenderId"`
// 	RecipientId string    `json:"RecipientId"`
// 	Message     string    `json:"Message"`
// 	CreatedAt   time.Time `json:"CreatedAt"`
// }

// type GroupMessage struct {
// 	Id        string `json:"Id"`
// 	SenderId  string `json:"SenderId"`
// 	GroupId   string `json:"GroupId"`
// 	Message   string `json:"Message"`
// 	CreatedAt string `json:"CreatedAt"`
// }

type PostFromClient struct {
	Id                  string    `json:"id"`
	Title               string    `json:"title"`
	Body                string    `json:"body"`
	CreatedAt           time.Time `json:"createdAt"`
	PrivacyLevel        string    `json:"privacyLevel"`
	CreatorId           string    `json:"creatorId"`
	Image               *Image    `json:"avatar,omitempty"`
	PostChosenFollowers []string  `json:"postChosenFollowers,omitempty"`
	GroupId             string    `json:"groupId,omitempty"`
}

// func (receivedData Post) parseForDB() *dbfuncs.Post {
// 	return &dbfuncs.Post{
// 		Title:        receivedData.Title,
// 		Body:         receivedData.Body,
// 		CreatorId:    receivedData.CreatorId,
// 		GroupId:      receivedData.GroupId,
// 		CreatedAt:    receivedData.CreatedAt,
// 		Image:        receivedData.Image.Data,
// 		PrivacyLevel: receivedData.PrivacyLevel,
// 	}
// }

type Notification struct {
	Id         string                 `json:"Id"`
	ReceiverId string                 `json:"ReceiverId"` // I changed this from RecieverId 4 Apr
	SenderId   string                 `json:"SenderId"`
	Payload    map[string]interface{} `json:"payload"`
	Body       string                 `json:"Body"`
	Type       string                 `json:"type"`
	CreatedAt  time.Time              `json:"CreatedAt"`
	Seen       bool                   `json:"Seen"`
}

// type NotificationSeen struct {
// 	Id string `json:"Id"`
// }

func (receivedData Notification) parseForDB(message string) *dbfuncs.Notification {
	return &dbfuncs.Notification{
		Body:       message,
		Type:       receivedData.Type,
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
	}
}

// func (dbNotification *dbfuncs.Notification) parseForClient() Notification {
// 	user := dbfunc.GetUserById(dbNotification.SenderId)
// 	body := map[string]interface{}{
// 		"message": dbNotification.Body,
// 		"user":    user,
// 	}
// 	return Notification{
// 		Type:       dbNotification.Type,
// 		ReceiverId: dbNotification.ReceiverId,
// 		SenderId:   dbNotification.SenderId,
// 	}
// }

// // these may not involve database calls but can still be sent through websockets
// // this can be resused for sending SignalReceiveds to other users about a user who has
// // e.g. registered, logged in, logged out, or changed their status
// type BasicUserInfo struct {
// 	Avatar string `json:"Avatar"`
// 	UserId         string `json:"UserId"`
// 	FirstName      string `json:"FirstName"`
// 	LastName       string `json:"LastName"`
// 	Nickname       string `json:"Nickname"`
// 	PrivacySetting string `json:"PrivacySetting"`
// }

// type RequestToFollow struct {
// 	User   BasicUserInfo `json:"User"`
// 	Status string        `json:"Status"`
// 	Type   string        `json:"Type"`
// }

type Follow struct {
	FollowerId  string `json:"FollowerId"`
	FollowingId string `json:"FollowingId"`
	Status      string `json:"Status"`
}

type AnswerRequestToFollow struct {
	SenderId   string `json:"SenderId"`
	ReceiverId string `json:"ReceiverId"`
	Reply      string `json:"Reply"`
}

type Unfollow struct {
	FollowerId  string `json:"FollowerId"`
	FollowingId string `json:"FollowingId"`
}

type Group struct {
	Id          string    `json:"Id"`
	CreatorId   string    `json:"CreatorId"`
	Title       string    `json:"Name"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

// type Event struct {
// 	EventId     string `json:"EventId"`
// 	GroupId     string `json:"GroupId"`
// 	Title       string `json:"Title"`
// 	Description string `json:"Description"`
// 	CreatorId   string `json:"CreatorId"`
// }

// func (receivedData Event) parseForDB() *dbfuncs.Event {
// 	return &dbfuncs.Event{
// 		EventId:     receivedData.EventId,
// 		GroupId:     receivedData.GroupId,
// 		Title:       receivedData.Title,
// 		Description: receivedData.Description,
// 		CreatorId:   receivedData.CreatorId,
// 	}
// }

// func (receivedData Comment) parseForDB() *dbfuncs.Comment {
// 	return &dbfuncs.Comment{
// 		Body:      receivedData.Body,
// 		CreatorId: receivedData.UserId,
// 		PostId:    receivedData.PostId,
// 	}
// }

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
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	cookie, err := r.Cookie("user_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	valid, err := dbfuncs.ValidateCookie(cookie.Value)
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
		fmt.Println(err, "defer")
		conn.Close()
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

	for id, arr := range activeConnections {
		fmt.Println(id, len(arr))
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

	// camelsBack:
	for {
		_, msgBytes, err := conn.ReadMessage()
		//possibly don't want to immediately delete the connection if there is an error
		if err != nil {
			fmt.Println("connections error not nil")
			myUpdatedConnections := []*websocket.Conn{}
			connectionLock.Lock()

			_, ok := activeConnections[userID]
			if !ok {
				connectionLock.Unlock()
				break
			}
			for _, c := range activeConnections[userID] {
				fmt.Println(userID, "userId")
				if c != conn {
					myUpdatedConnections = append(myUpdatedConnections, conn)
				}
				activeConnections[userID] = myUpdatedConnections
				if len(myUpdatedConnections) == 0 {
					delete(activeConnections, userID)
					// log.Println("User", userID, "disconnected, unable to read from websocket, error:", err)
				}
			}
			connectionLock.Unlock()
			break
		}

		err = broker(msgBytes, userID, conn, w)
		fmt.Println("returned from broker", err)
		if err != nil && err.Error() == "logout" {
			fmt.Println("logout error")
			break
		}

		// var finalStraw error
		// if err != nil {
		// 	finalStraw = notifyClientOfError(err, "error processing websocket message", userID)
		// }

		// if finalStraw != nil {Error sending user list to client: websocket: close sent
		// 	log.Println("error sending error message to client:", finalStraw)
		// 	break camelsBack
		// }

	}
}

// Check validity of signal and received data as far as possible. But get
// a working version first. We can add more checks later. Checks could
// also be called from the parseForDB methods.
func broker(msgBytes []byte, userID string, conn *websocket.Conn, w http.ResponseWriter) error {
	var signal SignalReceived
	err := json.Unmarshal(msgBytes, &signal)

	fmt.Println()
	if err != nil {
		log.Println("Error unmarshalling websocket signal:", err)
	}

	log.Println("signal.Type", signal.Type)

	switch signal.Type {
	// case "login":
	// No need. This is covered at the start of handleConnection.
	// Sign out and register:
	case "logout":
		logout(userID, conn, w)
		err = fmt.Errorf("logout")
		return err

	case "requestToFollow":
		var receivedData Follow
		unmarshalBody(signal.Body, &receivedData)
		fmt.Println(receivedData)
		err = requestToFollow(receivedData)
		fmt.Println("returned from requestToFollow")
	case "answerRequestToFollow":
		fmt.Println("case: answerRequestToFollow")
		var receivedData AnswerRequestToFollow
		unmarshalBody(signal.Body, &receivedData)
		log.Println("unmarshalled answer")
		err = answerRequestToFollow(receivedData)
		fmt.Println(err)
		fmt.Println("retured from answerRequestToFollow")
	case "unfollow":
		var receivedData Unfollow
		unmarshalBody(signal.Body, &receivedData)
		err = unfollow(receivedData)

	case "post":
		var receivedData PostFromClient
		unmarshalBody(signal.Body, &receivedData)
		err = post(receivedData)

	case "createGroup":
		var receivedData Group
		unmarshalBody(signal.Body, &receivedData)
		err = createGroup(receivedData)

	}

	log.Println("end of broker")
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
// as we don't have access to the cookie using websockets.
func logout(userID string, thisConn *websocket.Conn, w http.ResponseWriter) {
	fmt.Println("logging out")
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
	http.SetCookie(w, &http.Cookie{
		Name:     "user_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,

		SameSite: http.SameSiteLaxMode,
	})
}

func requestToFollow(receivedData Follow) error {
	fmt.Println("requestToFollow", receivedData)
	var follow dbfuncs.Follow
	follow.FollowingId = receivedData.FollowingId
	follow.FollowerId = receivedData.FollowerId

	log.Println(follow)

	private, err := dbfuncs.IsUserPrivate(receivedData.FollowingId)
	if err != nil {

		log.Println("error checking if user is public", err)
		return err

	}

	log.Println("private", private)

	if private {
		follow.Status = "pending"
	} else {
		follow.Status = "accepted"
	}

	err = dbfuncs.AddFollow(&follow)
	if err != nil {
		log.Println("error adding follow to database", err)
		notifyClientOfError(err, "requestToFollow", receivedData.FollowerId, nil)
		return err
	}

	follower, err := dbfuncs.GetUserById(receivedData.FollowerId)
	if err != nil {
		log.Println("error getting nickname from database", err)
		notifyClientOfError(err, "requestToFollow", receivedData.FollowerId, nil)
		return err
	}

	notification := Notification{
		ReceiverId: receivedData.FollowingId,
		SenderId:   receivedData.FollowerId,
		Type:       "notification requestToFollow",
	}

	data, _ := dbfuncs.GetUserById(receivedData.FollowerId)

	var message string
	if private {
		message = fmt.Sprintf("%s has requested to follow you", follower.Nickname)
	} else {
		message = fmt.Sprintf("%s has followed you", follower.Nickname)
	}

	notificationId, err := dbfuncs.AddNotification(notification.parseForDB(message))
	log.Println("notificationId", notificationId)
	if err != nil {
		log.Println("error adding requestToFollow notification to database", err)
		return err
	}
	notification.Id = notificationId
	log.Println("notification.Id", notification.Id)
	notification.Body = message
	notification.Payload = map[string]interface{}{
		"Message":   message,
		"Data":      data,
		"IsPrivate": private,
	}

	connectionLock.RLock()
	val, ok := activeConnections[follow.FollowingId]

	if ok {
		for _, c := range val {
			err = c.WriteJSON(notification)

			if err != nil {
				log.Println("error sending (potential) new follower info to recipient", err)
			}
		}
	}
	connectionLock.RUnlock()

	log.Println("err:", err)
	log.Println("receivedData.FollowerId", receivedData.FollowerId)
	notifyClientOfError(err, "requestToFollow", receivedData.FollowerId, nil)
	return err
}

// When received, client should request profile if they're on profile page.
func answerRequestToFollow(receivedData AnswerRequestToFollow) error {
	var err error

	log.Println("inside answerRequestToFollow")

	switch receivedData.Reply {
	case "yes":
		err = dbfuncs.AcceptFollow(receivedData.ReceiverId, receivedData.SenderId)
		if err != nil {
			log.Println("database error accepting follow", err)
			notifyClientOfError(err, "answerRequestToFollow accept", receivedData.SenderId, nil)
			return err
		}
	case "no":
		log.Println("case is no")
		log.Println(receivedData.ReceiverId, "ReceiverId\n", receivedData.SenderId, "SenderId")
		err := dbfuncs.DeleteFollow(receivedData.ReceiverId, receivedData.SenderId)
		if err != nil {
			log.Println("error rejecting follow", err)
			notifyClientOfError(err, "answerRequestToFollow reject", receivedData.SenderId, nil)
			return err
		}
	default:
		log.Println("unexpected reply in answerRequestToFollow:", receivedData.Reply)
		log.Printf("%s sent unexpected body %s, answering request from %s\n",
			receivedData.SenderId, receivedData.Reply, receivedData.ReceiverId)
		return fmt.Errorf("unexpected body in answerRequestToFollow")
	}

	notificationForDB := dbfuncs.Notification{
		Body:       receivedData.Reply,
		Type:       "answerRequestToFollow",
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
		Seen:       false,
	}

	notdificationId, err := dbfuncs.AddNotification(&notificationForDB)
	if err != nil {
		log.Println(err, "error adding notification to db")
	}

	notificationToSend := Notification{
		Id:         notdificationId,
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
		Type:       "notification answerRequestToFollow",
		CreatedAt:  notificationForDB.CreatedAt,
		Seen:       false,
	}

	following, err := dbfuncs.GetUserById(receivedData.SenderId)
	message := fmt.Sprintf("%s has said %s to your request to follow", following.Nickname, receivedData.Reply)
	notificationToSend.Body = message
	notificationToSend.Payload = map[string]interface{}{
		"Message": message,
		"Data":    receivedData,
	}

	connectionLock.RLock()
	for _, c := range activeConnections[receivedData.ReceiverId] {
		err = c.WriteJSON(notificationToSend)
		if err != nil {
			log.Println("error sending notification to recipient", err)
		}
	}
	connectionLock.RUnlock()

	fmt.Println(err)
	fmt.Println("success string")
	fmt.Println(receivedData.SenderId)
	// if receivedData.Reply == "yes" {
	// 	receiverFull, err := dbfuncs.GetUserById(receivedData.ReceiverId)
	// 	if err != nil {
	// 		log.Println(err, "error getting receiver by id")
	// 	}

	// 	receiver := BasicUserInfo{
	// 		Avatar:         receiverFull.Avatar,
	// 		Id:             receiverFull.Avatar,
	// 		FirstName:      receiverFull.FirstName,
	// 		LastName:       receiverFull.LastName,
	// 		Nickname:       receiverFull.Nickname,
	// 		PrivacySetting: receiverFull.PrivacySetting,
	// 	}

	// 	notifyClientOfError(err, "answerRequestToFollow "+receivedData.Reply, receivedData.SenderId, receiver)
	// 	return err
	// }
	notifyClientOfError(err, "answerRequestToFollow "+receivedData.Reply, receivedData.SenderId, receivedData.ReceiverId)
	log.Println("end of answer")
	return err
}

func unfollow(receivedData Unfollow) error {
	log.Println("reached unfollow")
	err := dbfuncs.DeleteFollow(receivedData.FollowerId, receivedData.FollowingId)
	if err != nil {
		log.Println("error deleting follow", err)
		notifyClientOfError(err, "unfollow", receivedData.FollowerId, nil)
		return err
	}

	notification := Notification{
		ReceiverId: receivedData.FollowingId,
		SenderId:   receivedData.FollowerId,
		Type:       "notification unfollow",
	}

	data, _ := dbfuncs.GetUserById(receivedData.FollowerId)
	message := fmt.Sprintf("%s has unfollowed you", data.Nickname)

	notificationId, err := dbfuncs.AddNotification(notification.parseForDB(message))
	log.Println("notificationId", notificationId)
	if err != nil {
		log.Println("error adding unfollow notification to database", err)
		return err
	}
	notification.Id = notificationId
	log.Println("unfollow notification.Id", notification.Id)
	notification.Body = message
	notification.Payload = map[string]interface{}{
		"Message": message,
		"Data":    data,
	}

	connectionLock.RLock()
	val, ok := activeConnections[receivedData.FollowingId]

	if ok {
		for _, c := range val {
			err = c.WriteJSON(notification)

			if err != nil {
				log.Println("error sending unfollow to following", err)
			}
		}
	}
	connectionLock.RUnlock()

	notifyClientOfError(err, "unfollow", receivedData.FollowerId, nil)
	return err
}

// Only notify a user of an error that occurred while processing an
// action they attempted. No need to notify someone if someone else
// failed to follow them, for example. I'm thinking, also, only notify
// user that a message couldn't be added to the db, since that affects
// them directly. If a message couldn't be sent to one of their connections,
// we can just log that and deal with it ourselves.
func notifyClientOfError(err error, message string, id string, whatever any) error {
	log.Println("notify client of error", err, message)

	data := map[string]interface{}{
		"message":  message,
		"whatever": whatever,
	}

	if err == nil {
		data["type"] = "success"
	} else {
		data["type"] = "error"
		data["error"] = err
	}

	connectionLock.RLock()

	val, ok := activeConnections[id]
	if ok {
		for _, c := range val {
			err = c.WriteJSON(data)
			if err != nil {
				fmt.Println("error sending error message to client:", err)
				break
			}
		}
	}
	connectionLock.RUnlock()
	return err
}

func createGroup(receivedData Group) error {
	log.Println("starting createGroup")
	group := dbfuncs.Group{
		CreatorId:   receivedData.CreatorId,
		Title:       receivedData.Title,
		Description: receivedData.Description,
	}
	err := dbfuncs.AddGroup(&group)
	if err != nil {
		log.Println("error adding group to database", err)
		notifyClientOfError(err, "error adding group to database", receivedData.CreatorId, nil)
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
		notifyClientOfError(err, "error adding group member to database", receivedData.CreatorId, nil)
		return err
	}

	notification := Notification{
		Type:     "createGroup",
		SenderId: receivedData.CreatorId,
		Seen:     false,
	}

	// We don't need to add this notification to the database. It's only relevant to users who are currently online and looking at the list of group.

	message := fmt.Sprintf("%s has created a new group, %s", group.CreatorId, group.Title)
	notification.Body = message
	notification.Payload = map[string]interface{}{
		"Message": message,
		"Data":    group,
	}

	connectionLock.RLock()
	for user := range activeConnections {
		for _, c := range activeConnections[user] {
			err = c.WriteJSON(notification)
			if err != nil {
				log.Println("error sending new group to client", err)
			}
		}
	}
	connectionLock.RUnlock()

	log.Println("err:", err)
	notifyClientOfError(err, "createGroup", receivedData.CreatorId, nil)
	return err
}

func post(receivedData PostFromClient) error {
	err := validateContent(receivedData.Body)
	if err != nil {
		return err
	}

	dbPost := dbfuncs.Post{
		Title:        receivedData.Title,
		Body:         receivedData.Body,
		CreatorId:    receivedData.CreatorId,
		PrivacyLevel: receivedData.PrivacyLevel,
	}

	if receivedData.Image != nil {
		dbPost.Image = receivedData.Image.Data
	}

	err = dbfuncs.AddPost(&dbPost)
	if err != nil {
		log.Println("error adding post to database", err)
		notifyClientOfError(err, "post", receivedData.CreatorId, nil)
		return err
	}

	receivedData.Id = dbPost.Id

	signalBody, err := json.Marshal(receivedData)
	if err != nil {
		log.Println("error marshalling receivedData", err)
		notifyClientOfError(err, "post", receivedData.CreatorId, nil)
		return err
	}

	signal := SignalReceived{
		Type: "post",
		Body: signalBody,
	}

	switch receivedData.PrivacyLevel {
	case "public":
		connectionLock.RLock()
		for user := range activeConnections {
			log.Println("user", user, "creatorId", receivedData.CreatorId)
			for _, c := range activeConnections[user] {
				fmt.Println(c, "c", signal, "signal")
				err = c.WriteJSON(signal)
				if err != nil {
					log.Println("error sending new post to clients", err)
				}
			}
		}
		connectionLock.RUnlock()
	case "private":
		for _, c := range activeConnections[receivedData.CreatorId] {
			err = c.WriteJSON(signal)
			if err != nil {
				log.Println("error sending new post to self", err)
			}
		}
		followers, err := dbfuncs.GetAcceptedFollowerIdsByFollowingId(receivedData.CreatorId)
		if err != nil {
			log.Println("error getting followers from database", err)
			notifyClientOfError(err, "post", receivedData.CreatorId, nil)
		}
		for _, followerId := range followers {
			for _, c := range activeConnections[followerId] {
				err = c.WriteJSON(signal)
				if err != nil {
					log.Println("error sending new post to client", err)
				}
			}
		}
	case "superprivate":
		for _, c := range activeConnections[receivedData.CreatorId] {
			err = c.WriteJSON(signal)
			if err != nil {
				log.Println("error sending new post to self", err)
			}
		}
		for _, followerId := range receivedData.PostChosenFollowers {
			postChosenFollower := dbfuncs.PostChosenFollower{
				PostId:     dbPost.Id,
				FollowerId: followerId,
			}
			err = dbfuncs.AddPostChosenFollower(&postChosenFollower)
			if err != nil {
				log.Println("error adding postChosenFollower to database", err)
				notifyClientOfError(err, "post", receivedData.CreatorId, nil)
				return err
			}
			for _, c := range activeConnections[followerId] {
				err = c.WriteJSON(signal)
				if err != nil {
					log.Println("error sending new post to client", err)
				}
			}
		}
	}

	notifyClientOfError(err, "post", receivedData.CreatorId, nil)
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
