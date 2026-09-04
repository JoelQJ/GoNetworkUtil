package tcp

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/JoelQJ/GoNetworkUtil/codec"
)

const DefaultMaxFrameSize = 10 * 1024 * 1024

type Client[T any] struct {
	Data *T

	conn     net.Conn
	order    binary.ByteOrder
	maxFrame int
	opts     *ConnectionOptions[T]

	once     sync.Once
	closeErr error
}

type Server[T any] struct {
	factory func(net.Conn) *Client[T]
}

type ConnectionOptions[T any] struct {
	ByteOrder    binary.ByteOrder
	MaxFrameSize int

	OnConnect    func(*Client[T])
	OnRawPacket  func(*Client[T], *codec.ByteBuf)
	OnDisconnect func(*Client[T])
}
