package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		fmt.Println("message type: ", messageType, "p: ", message, "err: ", err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("error reading message, unexpected closeError: ", err)
			}
			fmt.Println("error reading message, NOT unexpected closeError: ", err)
			return
		}
		fmt.Println("message recieved: ", string(message))
		if err := conn.WriteMessage(messageType, message); err != nil {
			fmt.Println("error writing message: ", err)
			return
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("unable to upgrade to websocket: ", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client Connected")
	err = conn.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		fmt.Println("writing message error: ", err)
	}
	reader(conn)
}
