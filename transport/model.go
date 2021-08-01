package transport

import "time"

type Head struct {
	Action   string    `json:"action"`
	From     string    `json:"from"`
	CreateAt time.Time `json:"create_at"`
}

type Payload struct {
	Head
	Data interface{} `json:"data"`
}
