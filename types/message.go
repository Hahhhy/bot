package types

type Message struct {
	Text     string
	GroupID  int64
	UserID   int64
	Nickname string
}

type APIRequest struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo,omitempty"`
}
