package codec

import "math"

func (self *ByteBuf) WriteBytes(data []byte){
	self.ensureCapacity(len(data))
	copy(
		self.buf[self.writeIndex:self.writeIndex+len(data)],
		data,
	)
	self.writeIndex += len(data)
}

func (self *ByteBuf) ReadBytes(data []byte){
	copy(
		data,
		self.buf[self.readIndex:self.readIndex+len(data)],
	)
	self.readIndex += len(data)
}

func (self *ByteBuf) WriteUInt8(integer uint8){
	self.ensureCapacity(SizeInt8)
	self.buf[self.writeIndex] = integer
	self.writeIndex += SizeInt8
}

func (self *ByteBuf) WriteInt8(integer int8){
	self.WriteUInt8(uint8(integer))
}

func (self *ByteBuf) ReadUInt8() uint8{
	integer := self.buf[self.readIndex]
	self.readIndex += SizeInt8
	return integer
}

func (self *ByteBuf) ReadInt8() int8{
	return int8(self.ReadUInt8())
}

func (self *ByteBuf) WriteUInt16(integer uint16){
	self.ensureCapacity(SizeInt16)
	self.order.PutUint16(self.buf[self.writeIndex:self.writeIndex+SizeInt16], integer)
	self.writeIndex += SizeInt16
}

func (self *ByteBuf) WriteInt16(integer int16){
	self.WriteUInt16(uint16(integer))
}

func (self *ByteBuf) ReadUInt16() uint16{
	integer := self.order.Uint16(self.buf[self.readIndex:self.readIndex+SizeInt16])
	self.readIndex += SizeInt16
	return integer
}

func (self *ByteBuf) ReadInt16() int16{
	return int16(self.ReadUInt16())
}

func (self *ByteBuf) WriteUInt32(integer uint32){
	self.ensureCapacity(SizeInt32)
	self.order.PutUint32(self.buf[self.writeIndex:self.writeIndex+SizeInt32], integer)
	self.writeIndex += SizeInt32
}

func (self *ByteBuf) WriteInt32(integer int32){
	self.WriteUInt32(uint32(integer))
}

func (self *ByteBuf) ReadUInt32() uint32{
	integer := self.order.Uint32(self.buf[self.readIndex:self.readIndex+SizeInt32])
	self.readIndex += SizeInt32
	return integer
}

func (self *ByteBuf) ReadInt32() int32{
	return int32(self.ReadUInt32())
}

func (self *ByteBuf) WriteUInt64(integer uint64){
	self.ensureCapacity(SizeInt64)
	self.order.PutUint64(self.buf[self.writeIndex:self.writeIndex+SizeInt64], integer)
	self.writeIndex += SizeInt64
}

func (self *ByteBuf) WriteInt64(integer int64){
	self.WriteUInt64(uint64(integer))
}

func (self *ByteBuf) ReadUInt64() uint64{
	integer := self.order.Uint64(self.buf[self.readIndex:self.readIndex+SizeInt64])
	self.readIndex += SizeInt64
	return integer
}

func (self *ByteBuf) ReadInt64() int64{
	return int64(self.ReadUInt64())
}

func (self *ByteBuf) WriteFloat32(float float32){
	self.WriteUInt32(math.Float32bits(float))
}

func (self *ByteBuf) ReadFloat32() float32{
	return math.Float32frombits(self.ReadUInt32())
}

func (self *ByteBuf) WriteFloat64(double float64){
	self.WriteUInt64(math.Float64bits(double))
}

func (self *ByteBuf) ReadFloat64() float64{
	return math.Float64frombits(self.ReadUInt64())
}