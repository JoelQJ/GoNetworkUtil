package packet

import (
	"fmt"

	"github.com/JoelQJ/GoNetworkUtil/codec"
)

type Dispatcher[T any] struct {
	handlers map[uint16]Handler[T]
}

func NewDispatcher[T any]() *Dispatcher[T] {
	return &Dispatcher[T]{handlers: make(map[uint16]Handler[T])}
}

func (d *Dispatcher[T]) Register(id uint16, handler Handler[T]) {
	if _, ok := d.handlers[id]; ok {
		panic(fmt.Sprintf("packet: handler with id %d already registered", id))
	}
	d.handlers[id] = handler
}

func (d *Dispatcher[T]) Dispatch(client T, buf *codec.ByteBuf) error {
	id := buf.ReadUInt16()

	handler, ok := d.handlers[id]
	if !ok {
		return fmt.Errorf("packet: no handler registered for id %d", id)
	}

	handler(client, buf)
	return nil
}