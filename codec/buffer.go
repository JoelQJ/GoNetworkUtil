package codec

import (
	"encoding/binary"
)

const (
	SizeBool    = 1
	SizeInt8    = 1
	SizeUInt8   = 1
	SizeByte    = 1
	SizeInt16   = 2
	SizeUInt16  = 2
	SizeInt32   = 4
	SizeUInt32  = 4
	SizeInt64   = 8
	SizeUInt64  = 8
	SizeFloat32 = 4
	SizeFloat64 = 8
	BytesToExpand = 100
)

type ByteBuf struct{
	buf []byte
	readIndex int
	writeIndex int
	order binary.ByteOrder
	MaxSize int
}

func New(order binary.ByteOrder) *ByteBuf{
	return &ByteBuf{order: order}
}

func Wrap(data []byte, order binary.ByteOrder) *ByteBuf{
	return &ByteBuf{buf: data, writeIndex: len(data), order: order}
}

func (self *ByteBuf) SetMaxSize(size int){
	self.MaxSize = size
}

func (self *ByteBuf) checkSize(size int){
	if self.MaxSize > 0 && size > self.MaxSize {
		panic("size exceeds max")
	}
}

func Write[T any](self *ByteBuf, codec Codec[T], value T){
	codec.Encode(self, value)
}
func Read[T any](self *ByteBuf, decodec Decodec[T]) T{
	return decodec.Decode(self)
}

func (self *ByteBuf) ensureCapacity(needed int){
	self.checkSize(needed)
	if self.writeReaming() < needed{
		var actualCapacity int = len(self.buf)
		var newCapacity = max(actualCapacity * 2, actualCapacity + needed + BytesToExpand)
		var newBuf []byte = make([]byte, newCapacity)
		copy(newBuf, self.buf[:self.writeIndex])
		self.buf = newBuf
	}
}

func (self *ByteBuf) writeReaming() int{
	return len(self.buf) - self.writeIndex
}

func (self *ByteBuf) ToBytesSlice() []byte{
	var maxPosition int = self.writeIndex;
	//Generamos un nuevo descriptor del slice pero apunta al mismo no se copia
	return self.buf[: maxPosition : maxPosition]
}

func (self *ByteBuf) Len() int{
	return self.writeIndex - self.readIndex
}

