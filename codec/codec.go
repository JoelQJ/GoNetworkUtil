package codec

type Codec[T any] interface{
	Encode(*ByteBuf, T)
	Decode(*ByteBuf) T
}