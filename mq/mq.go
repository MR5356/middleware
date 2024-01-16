package mq

type MessageQueue interface {
	Subscribe(topic string, callback interface{}) error
	Publish(topic string, data interface{})
}
