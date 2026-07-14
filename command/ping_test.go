package command

import (
	"testing"
)

func TestPingCommand(t *testing.T) {
	msg := Message{
		Text:     "/ping",
		GroupID:  123456,
		UserID:   1001,
		Nickname: "tester",
	}

	reply := Dispatch(msg)
	if reply == nil {
		t.Fatal("expected reply, got nil")
	}
	if reply.Action != "send_group_msg" {
		t.Errorf("expected send_group_msg, got %s", reply.Action)
	}
}
