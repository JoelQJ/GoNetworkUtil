package packet

import "GoNetworkUtils/codec"

type Packet interface {
	Id() uint32
	OnFinishDecode()
}


type PacketDecoder struct{
	codec.Decodec[Packet]
}

type PacketEncoder struct{
	codec.Codec[Packet]
}