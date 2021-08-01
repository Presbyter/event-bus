package nats

import (
	"github.com/Presbyter/evnet-bus/transport"
	"strings"
)

type Receiver struct {
	dsnArray []string
}

func (r *Receiver) QueueSubscribe(topic, queue string) (chan *transport.Payload, error) {
	nc, err := Connect(strings.Join(r.dsnArray, ","))
	if err != nil {
		return nil, err
	}

	subCh := make(chan *transport.Payload, 1<<6)

	if queue != "" {
		_, err = nc.QueueSubscribe(topic, queue, func(val *transport.Payload) {
			subCh <- val
		})
		if err != nil {
			return nil, err
		}
	} else {
		_, err = nc.Subscribe(topic, func(val *transport.Payload) {
			subCh <- val
		})
		if err != nil {
			return nil, err
		}
	}

	return subCh, nil
}

func NewReceiver(dsn ...string) *Receiver {
	return &Receiver{dsnArray: dsn}
}
