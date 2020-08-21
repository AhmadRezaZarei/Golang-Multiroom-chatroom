package main

import (
	"github.com/gorilla/websocket"
)

var hubGroup map[string]*Hub

func addNewHub(hubName string, hub *Hub) {


	if hubGroup==nil {
		hubGroup = make(map[string]*Hub)
	}
	hubGroup[hubName] = hub
	go hub.run()


}


func deleteHub(hubName string) {
	if hub, exists := hubGroup[hubName]; exists {

		for client, isConnected := range hub.clients {
			if isConnected {
				hub.unregister <- client
				delete(hub.clients, client)
				close(client.send)
			}
		}
		delete(hubGroup, hubName)
	}
}

func registerClient(hubName string, conn *websocket.Conn) {
	if hub , ok := hubGroup[hubName]; ok {
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		hub.register <- client
		client.hub.register <- client
		go client.writePump()
		go client.readPump()
		return
	}

	addNewHub(hubName, newHub())
	client := &Client{hub: hubGroup[hubName], conn: conn, send: make(chan []byte, 256)}
	hubGroup[hubName].register <- client
	go client.writePump()
	go client.readPump()
}
