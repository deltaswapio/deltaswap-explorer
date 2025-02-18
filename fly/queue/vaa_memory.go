package queue

import (
	"context"

	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

// VAAInMemoryOption represents a VAA queue in memory option function.
type VAAInMemoryOption func(*VAAInMemory)

// VAAInMemory represents VAA queue in memory.
type VAAInMemory struct {
	ch   chan Message
	size int
}

// NewVAAInMemory creates a VAA queue in memory instances.
func NewVAAInMemory(opts ...VAAInMemoryOption) *VAAInMemory {
	m := &VAAInMemory{size: 100}
	for _, opt := range opts {
		opt(m)
	}
	m.ch = make(chan Message, m.size)
	return m
}

// WithSize allows to specify an channel size when setting a value.
func WithSize(v int) VAAInMemoryOption {
	return func(i *VAAInMemory) {
		i.size = v
	}
}

// Publish sends the message to a channel.
func (i *VAAInMemory) Publish(_ context.Context, v *vaa.VAA, data []byte) error {
	i.ch <- &memoryConsumerMessage{
		data: data,
	}
	return nil
}

// Consume returns the channel with the received messages.
func (i *VAAInMemory) Consume(_ context.Context) <-chan Message {
	return i.ch
}

type memoryConsumerMessage struct {
	data []byte
}

func (m *memoryConsumerMessage) Data() []byte {
	return m.data
}

func (m *memoryConsumerMessage) Done(_ context.Context) {}

func (m *memoryConsumerMessage) Failed() {}

func (m *memoryConsumerMessage) IsExpired() bool {
	return false
}
