package packet

import "github.com/JoelQJ/GoNetworkUtil/codec"

type Packet any

type Encoder interface {
	ID() uint16
	Encode(*codec.ByteBuf)
}

type Decoder func(*codec.ByteBuf) Packet

type Handler[T any] func(T, Packet)