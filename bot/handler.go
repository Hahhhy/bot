package bot

import (
	"encoding/json"
	"log"
	"strings"

	"qqbot/command"
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

func HandleEvent(raw []byte) *command.APIRequest {
	var event Event
	if err := json.Unmarshal(raw, &event); err != nil {
		log.Printf("JSON parse error: %v, raw: %s", err, string(raw))
		return nil
	}

	if event.PostType == "message" && event.MessageType == "group" {
		log.Printf("Received group message from %s (UID:%d) in group %d",
			event.Sender.Nickname, event.Sender.UserID, event.GroupID)

		text := extractText(event.Message)
		if text == "" {
			return nil
		}

		msg := command.Message{
			Text:     text,
			GroupID:  event.GroupID,
			UserID:   event.Sender.UserID,
			Nickname: event.Sender.Nickname,
		}
		return command.Dispatch(msg)
	}
	return nil
}

type messageSegment struct {
	Type string `json:"type"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

func extractText(raw json.RawMessage) string {
	var segments []messageSegment
	if err := json.Unmarshal(raw, &segments); err != nil {
		return ""
	}
	var parts []string
	for _, seg := range segments {
		if seg.Type == "text" {
			parts = append(parts, seg.Data.Text)
		}
	}
	return strings.TrimSpace(strings.Join(parts, ""))
}
