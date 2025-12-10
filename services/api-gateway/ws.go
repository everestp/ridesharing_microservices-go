package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade  failed : %v", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("No  User ID provided")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error Reading  message : %v", err)
			break
		}
		log.Printf("Received message : %s", message)
	}

}
func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade  failed : %v", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("No  User ID provided")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Println("No  packageSlug ID provided")
		return
	}


	type Driver struct {
          ID  string  `json:"id"`
		  Name string  `json:"name"`
		  ProfilePicture string  `json:"profilePicture"`
		  CarPalte string  `json:"carPlate"`
		  PackageSlug string  `json:"packageSlug"`

	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID: userID,
			Name: "Everest",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPalte: "BHE4pa",
			PackageSlug: packageSlug,
		},
	}
	if err := conn.WriteJSON(msg); err != nil{
		log.Printf("Error sending message: %v",err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error Reading  message : %v", err)
			break
		}
		log.Printf("Received message : %s", message)
	}

}


