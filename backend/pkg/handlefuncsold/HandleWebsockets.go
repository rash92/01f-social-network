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
	Id         string    `json:"Id"`
	ReceiverId string    `json:"ReceiverId"` // I changed this from RecieverId 4 Apr
	SenderId   string    `json:"SenderId"`
	Body       string    `json:"Body"`
	Type       string    `json:"Type"`
	CreatedAt  time.Time `json:"CreatedAt"`
	Seen       bool      `json:"Seen"`
}

// type NotificationSeen struct {
// 	Id string `json:"Id"`
// }

func (receivedData Notification) parseForDB() *dbfuncs.Notification {
	return &dbfuncs.Notification{
		Body:       receivedData.Body,
		Type:       receivedData.Type,
		ReceiverId: receivedData.ReceiverId,
		SenderId:   receivedData.SenderId,
	}
}

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

		// if finalStraw != nil {
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

	log.Println(signal.Type)

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

		err = answerRequestToFollow(receivedData)
		fmt.Println(err)
		fmt.Println("retured from answerRequestToFollow")
	case "unfollow":
		var receivedData Unfollow
		unmarshalBody(signal.Body, &receivedData)
		err = unfollow(receivedData)
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
		notifyClientOfError(err, "requestToFollow", receivedData.FollowerId)
		return err
	}

	follower, err := dbfuncs.GetUserById(receivedData.FollowerId)
	if err != nil {
		log.Println("error getting nickname from database", err)
		notifyClientOfError(err, "requestToFollow", receivedData.FollowerId)
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
	notifyClientOfError(err, "requestToFollow", receivedData.FollowerId)
	return err
}

// When received, client should request profile if they're on profile page.
func answerRequestToFollow(receivedData AnswerRequestToFollow) error {
	var err error

	log.Println("inside requestToFollow")

	switch receivedData.Reply {
	case "yes":
		err = dbfuncs.AcceptFollow(receivedData.ReceiverId, receivedData.SenderId)
		if err != nil {
			log.Println("database error accepting follow", err)
			notifyClientOfError(err, "answerRequestToFollow accept", receivedData.SenderId)
			return err
		}
	case "no":
		log.Println("case is no")
		log.Println(receivedData.ReceiverId, "ReceiverId\n", receivedData.SenderId, "SenderId")
		err := dbfuncs.DeleteFollow(receivedData.ReceiverId, receivedData.SenderId)
		if err != nil {
			log.Println("error rejecting follow", err)
			notifyClientOfError(err, "answerRequestToFollow reject", receivedData.SenderId)
			return err
		}
	default:
		log.Println("unexpected reply in answerRequestToFollow:", receivedData.Reply)
		log.Printf("%s sent unexpected body %s, answering request from %s\n",
			receivedData.SenderId, receivedData.Reply, receivedData.ReceiverId)
		return fmt.Errorf("unexpected body in answerRequestToFollow")
	}

	// notificationForDB := dbfuncs.Notification{
	// 	Body:       receivedData.Reply,
	// 	Type:       "answerRequestToFollow",
	// 	CreatedAt:  time.Now(), // check format
	// 	ReceiverId: receivedData.SenderId,
	// 	SenderId:   receivedData.ReceiverId,
	// 	Seen:       false,
	// }

	// notificationToSend := Notification{
	// 	ReceiverId: receivedData.SenderId,
	// 	SenderId:   receivedData.ReceiverId,
	// 	Body:       receivedData.Reply,
	// 	Type:       "answerRequestToFollow",
	// 	CreatedAt:  notificationForDB.CreatedAt,
	// 	Seen:       false,
	// }

	// connectionLock.RLock()
	// for _, c := range activeConnections[receivedData.ReceiverId] {
	// 	err = c.WriteJSON(notificationToSend)
	// 	if err != nil {
	// 		log.Println("error sending notification to recipient", err)
	// 	}
	// }
	// connectionLock.RUnlock()

	// fmt.Println(err)
	// fmt.Println("success string")
	// fmt.Println(receivedData.SenderId)

	notifyClientOfError(err, "answerRequestToFollow", receivedData.SenderId)
	log.Println("end of answer")
	return err
}

func unfollow(receivedData Unfollow) error {
	err := dbfuncs.DeleteFollow(receivedData.FollowerId, receivedData.FollowingId)
	if err != nil {
		log.Println("error deleting follow", err)
		notifyClientOfError(err, "unfollow", receivedData.FollowerId)
		return err
	}

	notifyClientOfError(err, "unfollow", receivedData.FollowerId)
	return err
}

// Only notify a user of an error that occurred while processing an
// action they attempted. No need to notify someone if someone else
// failed to follow them, for example. I'm thinking, also, only notify
// user that a message couldn't be added to the db, since that affects
// them directly. If a message couldn't be sent to one of their connections,
// we can just log that and deal with it ourselves.
func notifyClientOfError(err error, message string, id string) error {
	log.Println("notify client of error", err, message)

	data := map[string]interface{}{
		"message": message,
	}

	if err != nil {
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
