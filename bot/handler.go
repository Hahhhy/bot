package bot

import (
	"encoding/json"
	"log"
	"strings"

	"qqbot/command"
	"qqbot/types"
)

type event struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	GroupID     int64  `json:"group_id"`
	Sender      struct {
		UserID   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
	} `json:"sender"`
	Message json.RawMessage `json:"message"`
}

func HandleEvent(raw []byte) *types.APIRequest {
	var e event
	if err := json.Unmarshal(raw, &e); err != nil {
		log.Printf("JSON parse error: %v, raw: %s", err, string(raw))
		return nil
	}

	if e.PostType == "message" && e.MessageType == "group" {
		log.Printf("Received group message from %s (UID:%d) in group %d",
			e.Sender.Nickname, e.Sender.UserID, e.GroupID)

		text := ExtractText(e.Message)
		if text == "" {
			return nil
		}

		msg := types.Message{
			Text:     text,
			GroupID:  e.GroupID,
			UserID:   e.Sender.UserID,
			Nickname: e.Sender.Nickname,
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

func ExtractText(raw json.RawMessage) string {
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
