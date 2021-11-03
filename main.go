package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

var ids []string

func send() {
	hubInstance := getInstance()
	fmt.Println("ids: ", ids)
	res0 := hubInstance.sendToClient(ids[0], Message{Type: 1, Body: "only to user 1"})
	if res0 {
		fmt.Println("Message 0 sent")
	}
	res1 := hubInstance.sendToClient(ids[1], Message{Type: 2, Body: "only to user 2"})
	if res1 {
		fmt.Println("Message 1 sent")
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
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
	hub := getInstance()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/send-example", sendExample)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
