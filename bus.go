package event_bus

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Presbyter/evnet-bus/transport"
	natsTrans "github.com/Presbyter/evnet-bus/transport/nats"
	"github.com/nats-io/nats.go"
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

func (b *bus) Send(event *transport.Payload) error {
	event.Head.From = b.cfg.AppName

	var err error
	topic := fmt.Sprintf("event-bus.%s.%s", b.cfg.AppName, event.Action)
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
	topic := fmt.Sprintf("event-bus.%s.>", b.cfg.AppName)
	queue := b.cfg.SubscribeQueue
	recCh, err := b.cfg.Transport.Subscribe(topic, queue)
	if err != nil {
		return err
	}

	go func() {
		for rec := range recCh {
			foo := b.cfg.SubscribeFunc
			if foo != nil {
				foo(rec)
			}
		}
	}()

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

type SubscribeFunc func(rec *transport.Payload)

var DefaultBusConfig = BusConfig{
	AppName:        "Default",
	SubscribeQueue: "2e71ad7b-7371-43fc-bc9a-6b69bc7c4bfb", // 使用uuid为了保证唯一性,独立的queue具有负载均衡效果
	SendRetryTimes: 0,
	Transport:      natsTrans.NewNatsTransport(nats.DefaultURL),
	SubscribeFunc: func(rec *transport.Payload) {
		_, _ = fmt.Fprintf(os.Stdout, "recevie message: %#v", *rec)
	},
	Log: os.Stdout,
}

func RegisterEventBus(cfg *BusConfig) *bus {
	return &bus{
		cfg: cfg,
	}
}
