package packet

import "github.com/JoelQJ/GoNetworkUtil/codec"

type Packet interface {
	Id() uint16
	OnFinishDecode()
	Encode(*codec.ByteBuf)
}

type Decoder interface {
	Decode(*codec.ByteBuf) Packet
}
