package packet

import (
	"fmt"

	"github.com/JoelQJ/GoNetworkUtil/codec"
)

type Dispatcher[T any] struct {
	decoders []Decoder
	handlers []Handler[T]
}

func NewDispatcher[T any](count uint16) *Dispatcher[T] {
	return &Dispatcher[T]{
		decoders: make([]Decoder, count),
		handlers: make([]Handler[T], count),
	}
}

func (d *Dispatcher[T]) RegisterDecoder(id uint16, decoder Decoder) {
	d.decoders[id] = decoder
}

func RegisterHandler[T any, P any](d *Dispatcher[T], id uint16, handler func(T, P)) {
	d.handlers[id] = func(client T, pkt Packet) {
		handler(client, pkt.(P))
	}
}

func (d *Dispatcher[T]) Decode(id uint16, buf *codec.ByteBuf) (Packet, error) {
	if int(id) >= len(d.decoders) || d.decoders[id] == nil {
		return nil, fmt.Errorf("no decoder registered for packet id %d", id)
	}
	return d.decoders[id](buf), nil
}

func (d *Dispatcher[T]) Dispatch(client T, id uint16, buf *codec.ByteBuf) error {
	pkt, err := d.Decode(id, buf)
	if err != nil {
		return err
	}
	if int(id) < len(d.handlers) && d.handlers[id] != nil {
		d.handlers[id](client, pkt)
	}
	return nil
}