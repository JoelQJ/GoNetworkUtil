package codec

func (b *ByteBuf) WriteStringUTF8(s string) {
	data := []byte(s)
	b.WriteUInt32(uint32(len(data)))
	b.WriteBytes(data)
}

func (b *ByteBuf) ReadStringUTF8() string {
	length := b.ReadUInt32()
	b.checkSize(int(length))
	data := make([]byte, length)
	b.ReadBytes(data)
	return string(data)
}
