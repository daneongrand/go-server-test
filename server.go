package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var clients map[*websocket.Conn]bool

type CommentList struct {
	Comment []Comment `json:"comments"`
}

type Comment struct {
	Text string `json:"text"`
}

var comments CommentList

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	clients = make(map[*websocket.Conn]bool)
	
	jsonFile, _ := os.Open("data.json")
	
	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	json.Unmarshal(byteValue, &comments)

	http.HandleFunc("/", wsHandler)
	http.HandleFunc("/comments", commentHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


func commentHandler (w http.ResponseWriter, r *http.Request) { 
	data := comments
	fmt.Println(data)
	tmpl, _ := template.ParseFiles("comments-temp.html")
	tmpl.Execute(w, data)
}


func wsHandler (w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()
	clients[connection] = true



	// for _, comment := range comments.Comment {
	// 	connection.WriteMessage(websocket.TextMessage, []byte(comment.Text))
	// }


	defer delete(clients, connection)
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		go addComment(message, &comments)
		go writeMessage(message)
		go loggerMessage(message)
	}

}




func rewriteData (commentList *CommentList) {
	jsonData, err := json.Marshal(&commentList)
	
	if err != nil {
		log.Println(err)
	}

	if err := ioutil.WriteFile("data.json", jsonData, 0); err != nil {
		log.Println(err)
	}
}


func loggerMessage(message []byte)  {
	log.Println(string(message))
}


func addComment(message []byte, commentList *CommentList) {

	comment := Comment{Text: string(message)}
	commentList.Comment = append(commentList.Comment, comment)
	fmt.Println(commentList)
	rewriteData(commentList)
}


func writeMessage (message []byte) {
	for connection := range clients {
		connection.WriteMessage(websocket.TextMessage, message)
	}
}





