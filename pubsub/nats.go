package pubsub

import (
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
	"time"
)

const DefaultURL =  "nats://39.105.141.168:4222"

var nc,_ = nats.Connect(DefaultURL,setupConnOptions([]nats.Option{})...)

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

// Send send msg to nc
func Send(subject string, msg []byte) error {
	err:= nc.Publish(subject,msg)
	if err != nil {
		fmt.Println(err)
	}
	nc.Flush()
	return nil
}

func Sub(subject string) {
	nc.Subscribe(subject,func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})


}