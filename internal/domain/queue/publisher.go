package queue

import "context"

type Publisher interface {
	Publish(context context.Context, body interface{})
}
