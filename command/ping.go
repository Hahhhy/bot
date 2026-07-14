package command

import "qqbot/types"

func init() {
	Register("/ping", PingHandler)
}

func PingHandler(msg types.Message) *types.APIRequest {
	return &types.APIRequest{
		Action: "send_group_msg",
		Params: map[string]interface{}{
			"group_id": msg.GroupID,
			"message": []map[string]interface{}{
				{"type": "text", "data": map[string]string{"text": "pong!"}},
			},
		},
	}
}
