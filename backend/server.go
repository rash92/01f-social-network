package main

import (
	"fmt"
	"log"
	"net/http"
	ws "social-network/backend/pkg/websockets"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

// func wsEndpoint(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "websocket endpoint")
// 	ws.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	ws, err = ws.
// }

func main() {
	fmt.Println("main func started")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", ws.WsEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
