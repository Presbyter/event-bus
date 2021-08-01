package nats

import (
	"github.com/nats-io/nats.go"
	"time"
)

var conn *nats.EncodedConn

// nats connect dsn: eg. nats://user:password@server:port

func Connect(dsn string) (*nats.EncodedConn, error) {
	if conn != nil {
		return conn, nil
	}

	options := []nats.Option{
		nats.Timeout(10 * time.Second),      // 设置超时时间
		nats.PingInterval(5 * time.Second),  // 设置ping间隔
		nats.MaxPingsOutstanding(3),         // 设置最大失败次数
		nats.ReconnectWait(2 * time.Second), // 重新连接
		nats.MaxReconnects(60),              // 尝试重新连接次数
	}

	nc, err := nats.Connect(dsn, options...)
	if err != nil {
		return nil, err
	}
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	conn = c
	return conn, nil
}

func ConnClose() {
	if conn == nil {
		return
	}

	conn.Close()
}
