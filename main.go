package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

// OneBot v11 接收的消息事件结构
type Event struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	GroupID     int64  `json:"group_id"`
	Sender      struct {
		UserID   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
	} `json:"sender"`
	Message json.RawMessage `json:"message"`
}

// OneBot v11 API 调用结构
type APIRequest struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo,omitempty"`
}

func main() {
	wsURL := "ws://localhost:3001"
	log.Printf("Connecting to NapCat at %s...", wsURL)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to NapCat!")

	// graceful shutdown on Ctrl+C
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

		var event Event
		if err := json.Unmarshal(raw, &event); err != nil {
			log.Printf("JSON parse error: %v, raw: %s", err, string(raw))
			continue
		}

		// 只处理群消息
		if event.PostType == "message" && event.MessageType == "group" {
			log.Printf("Received group message from %s (UID:%d) in group %d",
				event.Sender.Nickname, event.Sender.UserID, event.GroupID)

			reply := APIRequest{
				Action: "send_group_msg",
				Params: map[string]interface{}{
					"group_id": event.GroupID,
					"message": []map[string]interface{}{
						{"type": "text", "data": map[string]string{"text": "你好，我上线了！"}},
					},
				},
			}
			replyBytes, _ := json.Marshal(reply)
			if err := conn.WriteMessage(websocket.TextMessage, replyBytes); err != nil {
				log.Printf("Send error: %v", err)
			}
			log.Printf("Replied to group %d", event.GroupID)
		}
	}
}
