package event_bus

import (
	"fmt"
	"github.com/Presbyter/evnet-bus/transport"
	natsTrans "github.com/Presbyter/evnet-bus/transport/nats"
	"github.com/google/uuid"
	"strconv"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	t.Run("GenerateUUID", func(t *testing.T) {
		t.Log(uuid.New().String())
	})

	t.Run("EventBus", func(t *testing.T) {
		natsDsn := []string{
			"nats://127.0.0.1:4222",
		}

		cfg := DefaultBusConfig
		{
			cfg.BusName = "EventBusTesting"
			cfg.Transport = natsTrans.NewNatsTransport(natsDsn...)
			cfg.SubscribeQueue = "0e53731b-acc7-4dd3-a40d-bcb7098cd90e"
			cfg.SubscribeFunc = func(rec *transport.Payload) {
				switch rec.Action {
				case "test":
					fmt.Printf("action: test, msg: %#v\n", rec)
				default:
					fmt.Println(rec)
				}
			}
		}

		bus := RegisterEventBus(&cfg)

		t.Run("receive message", func(t *testing.T) {
			if err := bus.Receive(); err != nil {
				t.Errorf("receive message failed: %v", err)
			}
		})

		t.Run("send message", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				if err := bus.Send(&transport.Payload{
					Head: transport.Head{
						Action:   "test",
						CreateAt: time.Now(),
					},
					Data: "What are you kidding me?",
				}); err != nil {
					t.Errorf("send message failed: %v", err)
				}
			}
		})

		<-time.After(5 * time.Second)
	})

	t.Run("demo", func(t *testing.T) {
		ch := make(chan string, 0)

		go func() {
			for v := range ch {
				t.Log(v)
			}
		}()

		for i := 0; i < 100; i++ {
			ch <- strconv.Itoa(i)
		}

	})
}
