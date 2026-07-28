package codec

type Codec[T any] interface {
	Encode(*ByteBuf, T)
}

type Decodec[T any] interface {
	Decode(*ByteBuf) T
}
