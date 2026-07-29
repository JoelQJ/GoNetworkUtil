package packet

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"fmt"
)

type Dispatcher struct {
	decoders []Decoder
}

func NewDispatcher(count uint16) *Dispatcher {
	return &Dispatcher{
		decoders: make([]Decoder, count),
	}
}

func (d *Dispatcher) RegisterDecoder(id uint16, decoder Decoder) {
	d.decoders[id] = decoder
}

func (d *Dispatcher) Decode(id uint16, buf *codec.ByteBuf) (Packet, error) {
	if int(id) >= len(d.decoders) || d.decoders[id] == nil {
		return nil, fmt.Errorf("no decoder registered for packet id %d", id)
	}
	return d.decoders[id](buf), nil
}
