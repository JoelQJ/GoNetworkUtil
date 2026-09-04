package codec

import "math"

func (b *ByteBuf) WriteBytes(data []byte) {
	b.ensureCapacity(len(data))
	copy(b.buf[b.writeIndex:b.writeIndex+len(data)], data)
	b.writeIndex += len(data)
}

func (b *ByteBuf) ReadBytes(data []byte) {
	copy(data, b.buf[b.readIndex:b.readIndex+len(data)])
	b.readIndex += len(data)
}

func (b *ByteBuf) WriteSliceBytes(data []byte) {
	b.WriteInt32(int32(len(data)))
	b.WriteBytes(data)
}

func (b *ByteBuf) ReadSliceBytes() []byte {
	size := b.ReadInt32()
	b.checkSize(int(size))
	data := make([]byte, size)
	b.ReadBytes(data)
	return data
}

func (b *ByteBuf) WriteUInt8(v uint8) {
	b.ensureCapacity(SizeUInt8)
	b.buf[b.writeIndex] = v
	b.writeIndex += SizeUInt8
}

func (b *ByteBuf) WriteInt8(v int8) {
	b.WriteUInt8(uint8(v))
}

func (b *ByteBuf) WriteBool(v bool) {
	if v {
		b.WriteInt8(1)
	} else {
		b.WriteInt8(0)
	}
}

func (b *ByteBuf) ReadBool() bool {
	return b.ReadInt8() == 1
}

func (b *ByteBuf) ReadUInt8() uint8 {
	v := b.buf[b.readIndex]
	b.readIndex += SizeUInt8
	return v
}

func (b *ByteBuf) ReadInt8() int8 {
	return int8(b.ReadUInt8())
}

func (b *ByteBuf) WriteUInt16(v uint16) {
	b.ensureCapacity(SizeUInt16)
	b.order.PutUint16(b.buf[b.writeIndex:b.writeIndex+SizeUInt16], v)
	b.writeIndex += SizeUInt16
}

func (b *ByteBuf) WriteInt16(v int16) {
	b.WriteUInt16(uint16(v))
}

func (b *ByteBuf) ReadUInt16() uint16 {
	v := b.order.Uint16(b.buf[b.readIndex : b.readIndex+SizeUInt16])
	b.readIndex += SizeUInt16
	return v
}

func (b *ByteBuf) ReadInt16() int16 {
	return int16(b.ReadUInt16())
}

func (b *ByteBuf) WriteUInt32(v uint32) {
	b.ensureCapacity(SizeUInt32)
	b.order.PutUint32(b.buf[b.writeIndex:b.writeIndex+SizeUInt32], v)
	b.writeIndex += SizeUInt32
}

func (b *ByteBuf) WriteInt32(v int32) {
	b.WriteUInt32(uint32(v))
}

func (b *ByteBuf) ReadUInt32() uint32 {
	v := b.order.Uint32(b.buf[b.readIndex : b.readIndex+SizeUInt32])
	b.readIndex += SizeUInt32
	return v
}

func (b *ByteBuf) ReadInt32() int32 {
	return int32(b.ReadUInt32())
}

func (b *ByteBuf) WriteUInt64(v uint64) {
	b.ensureCapacity(SizeUInt64)
	b.order.PutUint64(b.buf[b.writeIndex:b.writeIndex+SizeUInt64], v)
	b.writeIndex += SizeUInt64
}

func (b *ByteBuf) WriteInt64(v int64) {
	b.WriteUInt64(uint64(v))
}

func (b *ByteBuf) ReadUInt64() uint64 {
	v := b.order.Uint64(b.buf[b.readIndex : b.readIndex+SizeUInt64])
	b.readIndex += SizeUInt64
	return v
}

func (b *ByteBuf) ReadInt64() int64 {
	return int64(b.ReadUInt64())
}

func (b *ByteBuf) WriteFloat32(v float32) {
	b.WriteUInt32(math.Float32bits(v))
}

func (b *ByteBuf) ReadFloat32() float32 {
	return math.Float32frombits(b.ReadUInt32())
}

func (b *ByteBuf) WriteFloat64(v float64) {
	b.WriteUInt64(math.Float64bits(v))
}

func (b *ByteBuf) ReadFloat64() float64 {
	return math.Float64frombits(b.ReadUInt64())
}
