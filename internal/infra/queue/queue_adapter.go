package queue

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/nats-io/nats.go"
)

type QueueAdapter struct {
	natsUrl   string
	conn      *nats.Conn
	listeners map[string][]Listener
}

func NewQueueAdapter(natsUrl string) *QueueAdapter {

	return &QueueAdapter{
		natsUrl:   natsUrl,
		listeners: make(map[string][]Listener),
	}
}

func (eb *QueueAdapter) ListenerRegister(eventType reflect.Type, handler func(*nats.Msg)) {
	eb.listeners[eventType.Name()] = append(eb.listeners[eventType.Name()], Listener{eventType, handler})
}

func (eb *QueueAdapter) Publish(ctx context.Context, eventPayload interface{}) error {
	eventType := reflect.TypeOf(eventPayload)
	payloadJson, _ := json.Marshal(eventPayload)

	log.Printf("--- Publish %s ---", eventType)

	for _, listener := range eb.listeners[eventType.Name()] {

		err := eb.conn.Publish(listener.eventType.Name(), []byte(payloadJson))

		if err != nil {
			return err
		}
	}

	return nil
}

func (eb *QueueAdapter) Connect(ctx context.Context) error {
	conn, err := nats.Connect(eb.natsUrl)
	if err != nil {
		log.Printf("Error connecting to nats --- %s", err.Error())
	}
	eb.conn = conn
	log.Printf("--- Connected to nats: %s ---", eb.conn.ConnectedAddr())
	return nil
}

func (eb *QueueAdapter) Disconnect(ctx context.Context) error {
	eb.conn.Close()
	log.Println("--- QueueAdapter disconnected ---")
	return nil
}

func (eb *QueueAdapter) StartConsuming(ctx context.Context, queueName string) error {
	var listeners []Listener
	for key, value := range eb.listeners {
		if key == queueName {
			listeners = value
		}
	}
	for _, v := range listeners {
		eb.conn.Subscribe(queueName, v.callback)
	}

	log.Printf("--- QueueAdapter StartConsuming queue %s ---", queueName)
	return nil
}
