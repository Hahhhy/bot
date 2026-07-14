package command

import "qqbot/types"

type Command func(msg types.Message) *types.APIRequest

var commands = map[string]Command{}

func Register(name string, cmd Command) {
	commands[name] = cmd
}

func Dispatch(msg types.Message) *types.APIRequest {
	for prefix, cmd := range commands {
		if len(msg.Text) >= len(prefix) && msg.Text[:len(prefix)] == prefix {
			return cmd(msg)
		}
	}
	return nil
}

func List() []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	return names
}
