package nats

import (
	"github.com/Presbyter/evnet-bus/transport"
)

type Transport struct {
	*Sender
	*Receiver
}

func (n Transport) Publish(topic string, msg *transport.Payload) error {
	if err := n.Sender.PublishMsg(topic, msg); err != nil {
		return err
	}

	return nil
}

func (n Transport) Subscribe(topic, queue string) (chan *transport.Payload, error) {
	msgCh, err := n.Receiver.QueueSubscribe(topic, queue)
	if err != nil {
		return nil, err
	}

	return msgCh, nil
}

func NewNatsTransport(dsn ...string) *Transport {
	return &Transport{
		Sender:   NewSender(dsn...),
		Receiver: NewReceiver(dsn...),
	}
}
