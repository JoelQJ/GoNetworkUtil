package udp

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/JoelQJ/GoNetworkUtil/codec"
)

type Client[T any] struct {
	Data *T

	conn  *net.UDPConn
	addr  *net.UDPAddr
	order binary.ByteOrder
	opts  *ConnectionOptions[T]

	once     sync.Once
	closeErr error
}

type Server[T any] struct {
	conn    *net.UDPConn
	factory func(*net.UDPAddr) *Client[T]
	opts    *ConnectionOptions[T]
}

type ConnectionOptions[T any] struct {
	ByteOrder binary.ByteOrder

	OnConnect   func(*Client[T])
	OnRawPacket func(*Client[T], *codec.ByteBuf)
}
