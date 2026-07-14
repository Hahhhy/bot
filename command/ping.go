package command

func init() {
	Register("/ping", PingHandler)
}

func PingHandler(msg Message) *APIRequest {
	return &APIRequest{
		Action: "send_group_msg",
		Params: map[string]interface{}{
			"group_id": msg.GroupID,
			"message": []map[string]interface{}{
				{"type": "text", "data": map[string]string{"text": "pong!"}},
			},
		},
	}
}
