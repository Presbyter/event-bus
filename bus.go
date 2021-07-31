package event_bus

import (
	"fmt"
	"github.com/Presbyter/evnet-bus/transport"
	"github.com/Presbyter/evnet-bus/transport/nats"
	"io"
	"os"
	"time"
)

type bus struct {
	cfg *BusConfig
}

func (b *bus) RegisterTransport(trans transport.AbstractTransport) *bus {
	if trans != nil {
		b.cfg.Transport = trans
	}
	return b
}

func (b *bus) Send(event *transport.SendPayload) error {
	var err error
	topic := fmt.Sprintf("event-bus/%s/%s", b.cfg.AppName, event.Action)
	if err = b.cfg.Transport.Publish(topic, event); err != nil {
		_, _ = fmt.Fprintf(b.cfg.Log, "event-bus publish message failed. error: %s", err)
	}

	retryTimes := b.cfg.SendRetryTimes
	if retryTimes <= 0 {
		return err
	}

	go func() {
		var success bool
		for i := 0; i < retryTimes; i++ {
			<-time.After(time.Second)
			_, _ = fmt.Fprintf(b.cfg.Log, "event-bus will retry publish message. current times: %d", i+1)
			if err = b.cfg.Transport.Publish(topic, event); err == nil {
				success = true
				break
			}
		}

		if !success {
			// 记录日志
			_, _ = fmt.Fprintf(b.cfg.Log, "[WARN] event-bus publish message failed. error: %s, event:%#v", err, *event)
		}
	}()

	return nil
}

func (b *bus) Receive() error {
	topic := fmt.Sprintf("event-bus/%s/*", b.cfg.AppName)
	rec, err := b.cfg.Transport.Subscribe(topic, b.cfg.SubscribeQueue)
	if err != nil {
		return err
	}

	foo := b.cfg.SubscribeFunc
	if foo != nil {
		return foo(rec)
	}

	return nil
}

type BusConfig struct {
	AppName        string
	SubscribeQueue string
	SendRetryTimes int
	Transport      transport.AbstractTransport
	SubscribeFunc  SubscribeFunc
	Log            io.Writer
}

type SubscribeFunc func(rec *transport.ReceivePayload) error

var DefaultBusConfig = BusConfig{
	AppName:        "Default",
	SubscribeQueue: "2e71ad7b-7371-43fc-bc9a-6b69bc7c4bfb", // 使用uuid为了保证唯一性,独立的queue具有负载均衡效果
	SendRetryTimes: 0,
	Transport:      &nats.NatsTrasport{},
	SubscribeFunc: func(rec *transport.ReceivePayload) error {
		_, _ = fmt.Fprintf(os.Stdout, "recevie message: %v", rec)
		return nil
	},
	Log: os.Stdout,
}

func RegisterEventBus(cfg *BusConfig) *bus {
	return &bus{
		cfg: cfg,
	}
}
