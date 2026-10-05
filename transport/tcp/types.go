package tcp

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/JoelQJ/GoNetworkUtil/packet"
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
	OnDisconnect func(*Client[T])

	Dispatcher *packet.Dispatcher[*Client[T]]
}
