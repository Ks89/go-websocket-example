package main

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool
	// Inbound messages from the clients.
	Broadcast chan []byte
	// Register requests from the clients.
	Register chan *Client
	// Unregister requests from clients.
	Unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) sendToClient(id string, msg Message) bool {
	fmt.Println("sendToClient called with id:: ", id)
	var c *Client
	for client := range h.Clients {
		if client.ID == id {
			fmt.Println("Client found: ", client)
			c = client
			if c != nil {
				c.Conn.WriteJSON(msg)
				return true
			}
			return false
		}
	}
	return false
}

func (h *Hub) run() {
	fmt.Println("run h: ", h)
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			fmt.Println("New client registered with ID: ", client.ID)
		case client := <-h.Unregister:
			fmt.Println("Unregistering client with ID: ", client.ID)
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			fmt.Println("Broadcasting message: ", message)
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
