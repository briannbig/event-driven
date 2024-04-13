package queue

type Publisher interface {
	Publish(body interface{}) error
}
