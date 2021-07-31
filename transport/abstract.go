package transport

type AbstractTransport interface {
	Publish(topic string, data *SendPayload) error
	Subscribe(topic, queue string) (*ReceivePayload, error)
}
