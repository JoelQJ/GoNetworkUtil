package udp

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/JoelQJ/GoNetworkUtil/packet"
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

	OnConnect func(*Client[T])

	Dispatcher *packet.Dispatcher[*Client[T]]
}
