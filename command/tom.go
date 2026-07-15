package command

import (
	"qqbot/image"
	"qqbot/types"
)

func init() {
	Register("/汤姆", TomHandler)
}

func TomHandler(msg types.Message) *types.APIRequest {
	path, err := image.GenerateTom(msg.Text)
	if err != nil {
		return nil
	}
	return &types.APIRequest{
		Action: "send_group_msg",
		Params: map[string]interface{}{
			"group_id": msg.GroupID,
			"message": []map[string]interface{}{
				{"type": "image", "data": map[string]string{"file": "file://" + path}},
			},
		},
	}
}
