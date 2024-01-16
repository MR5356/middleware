package mq

import "github.com/asaskevich/EventBus"

type MemoryMessageQueue struct {
	evbus EventBus.Bus
}

func NewMemoryMessageQueue() *MemoryMessageQueue {
	return &MemoryMessageQueue{
		evbus: EventBus.New(),
	}
}

func (m *MemoryMessageQueue) Subscribe(key string, callback interface{}) error {
	return m.evbus.Subscribe(key, callback)
}

func (m *MemoryMessageQueue) Publish(key string, data interface{}) {
	m.evbus.Publish(key, data)
}
