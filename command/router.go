package command

type Message struct {
	Text     string
	GroupID  int64
	UserID   int64
	Nickname string
}

type Command func(msg Message) *APIRequest

var commands = map[string]Command{}

func Register(name string, cmd Command) {
	commands[name] = cmd
}

func Dispatch(msg Message) *APIRequest {
	for prefix, cmd := range commands {
		if len(msg.Text) >= len(prefix) && msg.Text[:len(prefix)] == prefix {
			return cmd(msg)
		}
	}
	return nil
}

type APIRequest struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo,omitempty"`
}
