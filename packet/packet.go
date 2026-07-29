package packet

import "github.com/JoelQJ/GoNetworkUtil/codec"

type Packet interface {
	Id() uint16
	Encode(*codec.ByteBuf)
}

type Decoder func(*codec.ByteBuf) Packet
