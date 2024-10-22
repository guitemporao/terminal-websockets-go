package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	MessageType string
	Data        []byte
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	fmt.Printf("connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer conn.Close()

	send := make(chan Message)
	done := make(chan struct{})

	// Goroutine for reading messages
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("recv: %s\n", message)
		}
	}()

	// Goroutine for writing messages
	go func() {
		for {
			select {
			case msg := <-send:
				// Determine the message type
				var msgType int
				if msg.MessageType == "text" {
					msgType = websocket.TextMessage
				} else if msg.MessageType == "binary" {
					msgType = websocket.BinaryMessage
				}
				err := conn.WriteMessage(msgType, msg.Data)
				if err != nil {
					log.Println("write:", err)
				}
			case <-done:
				return
			}
		}
	}()

	// Read message from terminal
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter message (type 'text' or 'binary' followed by your message):...\n")

	for scanner.Scan() {
		text := scanner.Text()
		// Split the input to determine the message type
		var msgType string
		var msgContent []byte

		if len(text) > 5 && text[:5] == "text:" {
			msgType = "text"
			msgContent = []byte(text[5:])
		} else if len(text) > 7 && text[:7] == "binary:" {
			msgType = "binary"
			msgContent = []byte(text[7:])
		} else {
			fmt.Println("Please prefix your message with 'text:' or 'binary:'")
			continue
		}

		send <- Message{MessageType: msgType, Data: msgContent}
	}

	if err := scanner.Err(); err != nil {
		log.Println("read:", err)
	}
}
