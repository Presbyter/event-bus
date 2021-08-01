package transport

type AbstractTransport interface {
	Publish(topic string, data *Payload) error
	Subscribe(topic, queue string) (chan *Payload, error)
}
