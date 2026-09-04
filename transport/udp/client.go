package udp

import (
	"encoding/binary"
	"net"

	"github.com/JoelQJ/GoNetworkUtil/codec"
	"github.com/JoelQJ/GoNetworkUtil/packet"
)

func NewClient[T any](conn *net.UDPConn, addr *net.UDPAddr, data *T, opts *ConnectionOptions[T]) *Client[T] {
	if data == nil {
		data = new(T)
	}
	if opts == nil {
		opts = &ConnectionOptions[T]{}
	}
	order := opts.ByteOrder
	if order == nil {
		order = binary.BigEndian
	}
	return &Client[T]{
		Data:  data,
		conn:  conn,
		addr:  addr,
		order: order,
		opts:  opts,
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
	return NewClient(conn, addr, data, opts), nil
}

func (c *Client[T]) ReadPacket() (*codec.ByteBuf, error) {
	buf := make([]byte, 65535)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return codec.Wrap(buf[:n], c.order), nil
}

func (c *Client[T]) Send(encoder packet.Encoder) error {
	body := codec.New(c.order)
	body.WriteUInt16(encoder.ID())
	encoder.Encode(body)

	return c.SendRaw(body.Bytes())
}

func (c *Client[T]) SendRaw(data []byte) error {
	_, err := c.conn.WriteToUDP(data, c.addr)
	return err
}

func (c *Client[T]) RemoteAddr() net.Addr {
	return c.addr
}

func (c *Client[T]) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Client[T]) Close() error {
	c.once.Do(func() {
		c.closeErr = c.conn.Close()
	})
	return c.closeErr
}

func (c *Client[T]) Err() error {
	return c.closeErr
}
