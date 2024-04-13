package broker

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"prices/internal/config"
)

type MessageBrokerHandler func(msg *nats.Msg)

type MessageBroker interface {
	Consume(handler MessageBrokerHandler) error
	Publish(msg *nats.Msg) error
	Close()
}

type messageBroker struct {
	conn *nats.Conn
	cfg  *config.Config
}

func New(cfg *config.Config) (MessageBroker, error) {
	fmt.Println(cfg.Nats.URL)
	conn, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to nats %v", err)
	}

	return &messageBroker{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (m *messageBroker) Consume(handler MessageBrokerHandler) error {
	errChan := make(chan error)
	go m.subscribe(m.cfg.Nats.Queues.Validation, handler, errChan)

	for err := range errChan {
		return err
	}

	return nil
}

func (m *messageBroker) subscribe(channel string, handler MessageBrokerHandler, errChan chan error) {
	buffer := 64
	ch := make(chan *nats.Msg, buffer)
	sub, err := m.conn.ChanSubscribe(channel, ch)
	if err != nil {
		errChan <- err
		close(errChan)
		close(ch)
	}
	defer sub.Unsubscribe()

	for msg := range ch {
		handler(msg)
	}
}

func (m *messageBroker) Publish(msg *nats.Msg) error {
	return m.conn.PublishMsg(msg)
}

func (m *messageBroker) Close() {
	m.conn.Close()
}
