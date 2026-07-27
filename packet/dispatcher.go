package packet

import "GoNetworkUtils/codec"


type Dispatcher struct {
    decodecs []PacketDecoder
}


func (self *Dispatcher) Register(index uint16, decodec PacketDecoder){
	self.decodecs[index] = decodec
}

func (self *Dispatcher) Distpach(buf *codec.ByteBuf){
	var id uint16 = buf.ReadUInt16();
	var decoder PacketDecoder = self.decodecs[id]
	var packet = decoder.Decode(buf)
	packet.OnFinishDecode()
}