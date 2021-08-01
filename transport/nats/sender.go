package nats

import (
	"github.com/Presbyter/evnet-bus/transport"
	"strings"
)

type Sender struct {
	dsnArray []string
}

func (s *Sender) PublishMsg(topic string, msg *transport.Payload) error {
	nc, err := Connect(strings.Join(s.dsnArray, ","))
	if err != nil {
		return err
	}

	if err = nc.Publish(topic, msg); err != nil {
		return err
	}
	return nil
}

func NewSender(dsn ...string) *Sender {
	if dsn == nil || len(dsn) == 0 {
		panic("dsn array must not empty")
	}
	return &Sender{dsnArray: dsn}
}
