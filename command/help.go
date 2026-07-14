package command

import (
	"fmt"
	"strings"

	"qqbot/types"
)

func init() {
	Register("/帮助", HelpHandler)
	Register("/help", HelpHandler)
}

func HelpHandler(msg types.Message) *types.APIRequest {
	names := List()
	return &types.APIRequest{
		Action: "send_group_msg",
		Params: map[string]interface{}{
			"group_id": msg.GroupID,
			"message": []map[string]interface{}{
				{"type": "text", "data": map[string]string{"text": fmt.Sprintf("可用命令: %s", strings.Join(names, ", "))}},
			},
		},
	}
}
