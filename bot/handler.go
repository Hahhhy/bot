package bot

import (
	"encoding/json"
	"log"
)

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

type APIRequest struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo,omitempty"`
}

func HandleEvent(raw []byte) *APIRequest {
	var event Event
	if err := json.Unmarshal(raw, &event); err != nil {
		log.Printf("JSON parse error: %v, raw: %s", err, string(raw))
		return nil
	}

	if event.PostType == "message" && event.MessageType == "group" {
		log.Printf("Received group message from %s (UID:%d) in group %d",
			event.Sender.Nickname, event.Sender.UserID, event.GroupID)

		return &APIRequest{
			Action: "send_group_msg",
			Params: map[string]interface{}{
				"group_id": event.GroupID,
				"message": []map[string]interface{}{
					{"type": "text", "data": map[string]string{"text": "你好，我上线了！"}},
				},
			},
		}
	}
	return nil
}
