package codec

import "encoding/binary"

type ByteBuf struct {
	buf        []byte
	readIndex  int
	writeIndex int
	order      binary.ByteOrder
	maxSize    int
}

type Encoder[T any] interface {
	Encode(*ByteBuf, T)
}

type Decoder[T any] interface {
	Decode(*ByteBuf) T
}
