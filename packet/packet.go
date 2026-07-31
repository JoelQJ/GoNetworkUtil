package packet

import "github.com/JoelQJ/GoNetworkUtil/codec"

type Packet any

type Encoder interface {
	Id() uint16
	Encode(*codec.ByteBuf)
}

type Decoder func(*codec.ByteBuf) Packet
