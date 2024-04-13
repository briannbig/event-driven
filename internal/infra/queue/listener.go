package queue

import (
	"reflect"

	"github.com/nats-io/nats.go"
)

type Listener struct {
	eventType reflect.Type
	callback  func(*nats.Msg)
}
