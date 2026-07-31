package udp

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"github.com/JoelQJ/GoNetworkUtil/packet"
	"encoding/binary"
	"net"
)

type Client[T any] struct {
	Data      *T
	Alive     bool
	conn      *net.UDPConn
	addr      *net.UDPAddr
	byteOrder binary.ByteOrder
}

type ConnectionOptions[T any] struct {
	ByteOrder   binary.ByteOrder
	Dispatcher  *packet.Dispatcher
	OnConnect   func(*Client[T])
	OnRawPacket func(*Client[T], *codec.ByteBuf)
	OnPacket    func(*Client[T], packet.Packet)
}

func NewClient[T any](conn *net.UDPConn, addr *net.UDPAddr, data *T) *Client[T] {
	return &Client[T]{
		Data:  data,
		Alive: true,
		conn:  conn,
		addr:  addr,
	}
}

func Connect[T any](address string, data *T, opts *ConnectionOptions[T]) (*Client[T], error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	client := NewClient(conn, addr, data)
	ApplyOptions(client, opts)
	return client, nil
}

func ApplyOptions[T any](c *Client[T], opts *ConnectionOptions[T]) {
	if opts == nil {
		return
	}
	if opts.ByteOrder != nil {
		c.byteOrder = opts.ByteOrder
	} else {
		c.byteOrder = binary.BigEndian
	}
}

func (c *Client[T]) ReadRaw() (*codec.ByteBuf, error) {
	buf := make([]byte, 65535)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return codec.Wrap(buf[:n], c.byteOrder), nil
}

type PacketHandler[T any] interface {
	OnFinishDecode(*Client[T])
}

func (c *Client[T]) Send(encoder packet.Encoder) error {
	body := codec.New(c.byteOrder)
	body.WriteUInt16(encoder.Id())
	encoder.Encode(body)

	_, err := c.conn.Write(body.ToBytesSlice())
	return err
}

func (c *Client[T]) SendRaw(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *Client[T]) RemoteAddr() net.Addr {
	return c.addr
}

func (c *Client[T]) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Client[T]) Close() error {
	c.Alive = false
	return c.conn.Close()
}
