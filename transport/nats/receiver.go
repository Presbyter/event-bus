package nats

import "github.com/Presbyter/evnet-bus/transport"

type NatsTrasport struct{}

func (n NatsTrasport) Publish(topic string, data *transport.SendPayload) error {
	panic("implement me")
}

func (n NatsTrasport) Subscribe(topic, queue string) (*transport.ReceivePayload, error) {
	panic("implement me")
}
