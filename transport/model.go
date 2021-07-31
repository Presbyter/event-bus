package transport

import "time"

type Head struct {
	Action   string    `json:"action"`
	From     string    `json:"from"`
	CreateAt time.Time `json:"create_at"`
}

type SendPayload struct {
	Head
	Data interface{} `json:"data"`
}

type ReceivePayload struct {
	Head
	ReceiveAt time.Time   `json:"receive_at"`
	Data      interface{} `json:"data"`
}
