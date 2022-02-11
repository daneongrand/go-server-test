package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var clients map[*websocket.Conn]bool

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	clients = make(map[*websocket.Conn]bool)
	http.HandleFunc("/", swHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func swHandler (w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()
	clients[connection] = true
	defer delete(clients, connection)
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		go writeMessage(message)
		go loggerMessage(message)
	}

}

func loggerMessage(message []byte)  {
	log.Println(string(message))
}

func writeMessage (message []byte) {
	for connection := range clients {
		connection.WriteMessage(websocket.TextMessage, message)
	}
}