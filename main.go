package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)


// use postman
// 1. connect to the websocket to localhost:8080/ws
// 2. get request to  localhost:8080/send-example
// 3. see websocket output in postman. You should be able to see the message send from the server via send method

var ids []string

func send() {
	hubInstance := getInstance()
	fmt.Println("ids: ", ids)
	if len(ids) != 0 {
		res0 := hubInstance.sendToClient(ids[0], Message{Type: 1, Body: "only to user 1"})
		if res0 {
			fmt.Println("Message 0 sent")
		}
	}
	if len(ids) > 1 {
		res1 := hubInstance.sendToClient(ids[1], Message{Type: 2, Body: "only to user 2"})
		if res1 {
			fmt.Println("Message 1 sent")
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub := getInstance()
	// create a new client with data received from HTTP request
	// in this case to do a basic example I generate a random UUID
	// to identify this client
	newId := uuid.NewString()
	ids = append(ids, newId)
	client := &Client{
		ID: newId,
		Hub: hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	// register the new client to the hub
	client.Hub.Register <- client
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func sendExample(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	fmt.Println("sending")
	go send()
}

func main() {
	hubInstance := getInstance()
	go hubInstance.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/send-example", sendExample)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
