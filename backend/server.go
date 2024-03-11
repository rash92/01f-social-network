package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"server/pkg/db/dbfuncs"
	"server/pkg/handlefuncs"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
}

func DeleteUserByUsername(username string) error {
	stmt, err := db.Prepare("DELETE FROM Users WHERE firstname= ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}

// var (
// 	upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			origin := r.Header.Get("Origin")
// 			// Modify this to match your allowed origins
// 			return origin == "http://localhost:8000"

// 		},
// 	}
// 	activeConnections = make(map[*websocket.Conn]string)
// 	connectionLock    sync.Mutex
// 	messageLock       sync.Mutex

// 	// userListLock      sync.Mutex
// )

// func handleConnection(w http.ResponseWriter, r *http.Request) {
// 	cookie, err := r.Cookie("user_token")
// 	if err != nil || !dbfuncs.ValidateCookie(cookie.Value) {

// 		// If the cookie is not valid, close the WebSocket connection
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Error upgrading to WebSocket:", err)
// 		return
// 	}
// 	defer func() {
// 		// fmt.Println("conncetion closing")
// 		conn.Close()

// 	}()

// 	userID, err := dbfuncs.GetUserIdFromCokie(cookie.Value)
// 	if err != nil {
// 		fmt.Println("Error sending user ID to client:", err)
// 	}

// 	connectionLock.Lock()
// 	activeConnections[conn] = userID
// 	connectionLock.Unlock()
// 	isLoggingOut := false
// 	broadcastUserList(userID)
// 	for {

// 		_, p, err := conn.ReadMessage()
// 		if err != nil {
// 			connectionLock.Lock()
// 			if isLoggingOut {
// 				for client, id := range activeConnections {
// 					if id == userID && client != conn {
// 						fmt.Println(id, client)
// 						data := map[string]interface{}{
// 							"data": "log me out",
// 							"type": "logout",
// 						}
// 						err := client.WriteJSON(data)
// 						if err != nil {
// 							fmt.Println("Error logout message  to client:", err)
// 						}
// 					}
// 				}

// 			}
// 			delete(activeConnections, conn)
// 			connectionLock.Unlock()
// 			broadcastUserList(userID)
// 			fmt.Println("User", userID, "disconnected")

// 			break
// 		}

// 		var receivedData handlefuncs.Message

// 		err = json.Unmarshal(p, &receivedData)
// 		if receivedData.Type == "logout" {
// 			isLoggingOut = true
// 		}
// 		if err != nil {
// 			fmt.Println(err, "eror of  Unmarshal")
// 		}
// 		fmt.Println(receivedData, "receivedData message")
// 		messageLock.Lock()
// 		id, created, err := dbfuncs.AddMessage(receivedData.SenderID, receivedData.RecipientID, receivedData.Message, receivedData.Type)
// 		messageLock.Unlock()
// 		if err != nil {
// 			fmt.Println(err, "error adding message in the data base main line 120")
// 		}

// 		receivedData.ID = id.String()
// 		receivedData.Created = created.Format(time.RFC3339)
// 		message := map[string]interface{}{
// 			"data": receivedData,
// 			"type": receivedData.Type,
// 		}

// 		for client, id := range activeConnections {
// 			if id == receivedData.RecipientID {
// 				err := client.WriteJSON(message)
// 				if err != nil {
// 					fmt.Println("Error sending user list to client:", err)
// 				}

// 			}

// 		}
// 	}
// }

// func broadcastUserList(id string) {

// 	var userList []string

// 	connectionLock.Lock()

// 	for _, userID := range activeConnections {

// 		userList = append(userList, userID)

// 	}
// 	connectionLock.Unlock()
// 	// Broadcast the list to all connected clients

// 	message := map[string]interface{}{
// 		"data": userList,
// 		"type": "online-user",
// 	}
// 	for client := range activeConnections {

// 		err := client.WriteJSON(message)
// 		if err != nil {
// 			fmt.Println("Error sending user list to client:", err)
// 		}

// 	}

// }

func wrapperHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlefuncs.Cors(&w, r)
		cookie, err := r.Cookie("user_token")

		if err != nil || !dbfuncs.ValidateCookie(cookie.Value) {
			// fmt.Println("cookie value wrapper function", cookie.Value)
			http.Error(w, `{"error": "something went to wrong"}`, http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func main() {
	defer db.Close()
	handlefuncs.SetDatabase(db)
	dbfuncs.SetDatabase(db)
	// DeleteUserByUsername("bilal")
	// magarate.Magarate()
	// http.HandleFunc("/ws", handleConnection)

	http.HandleFunc("/newUser", handlefuncs.HandleNewUser)
	http.HandleFunc("/check-nickname", handlefuncs.HanndleUserNameIsDbOrEmail)
	http.HandleFunc("/check-email", handlefuncs.HanndleUserNameIsDbOrEmail)
	http.HandleFunc("/login", handlefuncs.HandleLogin)
	http.HandleFunc("/checksession", handlefuncs.HandleValidateCookie)
	http.HandleFunc("/add-post", handlefuncs.HandleAddPost)
	http.HandleFunc("/get-catogries", handlefuncs.HandleCatogries)
	http.HandleFunc("/get-posts", wrapperHandler(handlefuncs.HandleGetPosts))
	// Commented out because definition change to placeholder for the sake of the
	// web sockets.
	// http.HandleFunc("/add-Comment", handlefuncs.HandleAddComment)
	http.HandleFunc("/logout", handlefuncs.HandleLogout)
	http.HandleFunc("/react-Post-like-dislike", handlefuncs.HandlePostLikeDislike)
	// http.HandleFunc("/react-comment-like-dislike", handlefuncs.HandleCommenttLikeDislike)
	http.HandleFunc("/removePost", handlefuncs.HandleRemovePost)
	http.HandleFunc("/profile", wrapperHandler(handlefuncs.HandleGetProfile))
	http.HandleFunc("/get-users", handlefuncs.HandleGetUsers)
	http.HandleFunc("/get-messages", handlefuncs.MessagesHandler)
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./pkg/db/images"))))

	fmt.Println("Starting server on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
