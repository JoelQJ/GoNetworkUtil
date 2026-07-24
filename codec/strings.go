package codec

func (self *ByteBuf) WriteStringUTF8(string string){
	data := []byte(string)
	self.WriteUInt32(uint32(len(data)))
	self.WriteBytes(data)
}

func (self *ByteBuf) ReadStringUTF8() string{
	length := self.ReadUInt32()
	data := make([]byte, length)
	self.ReadBytes(data)
	return string(data)
}