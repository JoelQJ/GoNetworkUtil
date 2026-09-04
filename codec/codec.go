package codec

import "encoding/binary"

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

	defaultExpand = 100
)

func New(order binary.ByteOrder) *ByteBuf {
	return &ByteBuf{order: order}
}

func Wrap(data []byte, order binary.ByteOrder) *ByteBuf {
	return &ByteBuf{buf: data, writeIndex: len(data), order: order}
}

func (b *ByteBuf) SetMaxSize(size int) {
	b.maxSize = size
}

func Encode[T any](b *ByteBuf, enc Encoder[T], value T) {
	enc.Encode(b, value)
}

func Decode[T any](b *ByteBuf, dec Decoder[T]) T {
	return dec.Decode(b)
}

func (b *ByteBuf) checkSize(size int) {
	if b.maxSize > 0 && size > b.maxSize {
		panic("size exceeds max")
	}
}

func (b *ByteBuf) ensureCapacity(needed int) {
	b.checkSize(needed)
	if b.writeRemaining() < needed {
		actual := len(b.buf)
		newCap := max(actual*2, actual+needed+defaultExpand)
		newBuf := make([]byte, newCap)
		copy(newBuf, b.buf[:b.writeIndex])
		b.buf = newBuf
	}
}

func (b *ByteBuf) writeRemaining() int {
	return len(b.buf) - b.writeIndex
}

func (b *ByteBuf) Bytes() []byte {
	n := b.writeIndex
	return b.buf[:n:n]
}

func (b *ByteBuf) Len() int {
	return b.writeIndex - b.readIndex
}
