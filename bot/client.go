package bot

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

const wsURL = "ws://localhost:3001"

func Run() {
	log.Printf("Connecting to NapCat at %s...", wsURL)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to NapCat!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		log.Println("Shutting down...")
		conn.Close()
		os.Exit(0)
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		req := HandleEvent(raw)
		if req != nil {
			reply, _ := json.Marshal(req)
			if err := conn.WriteMessage(websocket.TextMessage, reply); err != nil {
				log.Printf("Send error: %v", err)
			}
			log.Printf("Replied to group")
		}
	}
}
