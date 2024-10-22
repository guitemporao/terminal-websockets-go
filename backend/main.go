package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	//1. upgrade to incoming GET request into websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket upgrade failed: ", err)
		return
	}

	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message failed: ", err)
			return
		}

		fmt.Printf("received message: %s\n", message)

		res := fmt.Sprintf("echo: %s", message, time.Now().Format(time.RFC3339))

		if err := conn.WriteMessage(messageType, []byte(res)); err != nil {
			fmt.Println("write message failed: ", err)
			return
		}
	}
}

func main() {
	// init websocket server
	http.HandleFunc("/ws", handler)
	fmt.Println("websocket server started", ":8080")

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("websocket server failed: ", err)
	}
}
