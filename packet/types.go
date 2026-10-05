package packet

import "github.com/JoelQJ/GoNetworkUtil/codec"

// Handler receives the connection the packet came from and the buffer already
// positioned right after the packet id.
type Handler[T any] func(T, *codec.ByteBuf)