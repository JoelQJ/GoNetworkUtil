package tcp

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/JoelQJ/GoNetworkUtil/codec"
	"github.com/JoelQJ/GoNetworkUtil/packet"
)

func NewClient[T any](conn net.Conn, data *T, opts *ConnectionOptions[T]) *Client[T] {
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
	maxFrame := opts.MaxFrameSize
	if maxFrame <= 0 {
		maxFrame = DefaultMaxFrameSize
	}
	return &Client[T]{
		Data:     data,
		conn:     conn,
		order:    order,
		maxFrame: maxFrame,
		opts:     opts,
	}
}

func Connect[T any](address string, data *T, opts *ConnectionOptions[T]) (*Client[T], error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, data, opts), nil
}

func (c *Client[T]) ReadLoop() {
	if c.opts.OnConnect != nil {
		c.opts.OnConnect(c)
	}
	defer func() {
		c.Close()
		if c.opts.OnDisconnect != nil {
			c.opts.OnDisconnect(c)
		}
	}()

	for {
		buf, err := c.ReadPacket()
		if err != nil {
			return
		}
		if c.opts.OnRawPacket != nil {
			c.opts.OnRawPacket(c, buf)
		}
	}
}

func (c *Client[T]) ReadPacket() (*codec.ByteBuf, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return nil, err
	}

	size := int32(codec.Wrap(header, c.order).ReadInt32())
	if size < 2 || int(size) > c.maxFrame {
		return nil, errors.New("invalid frame size")
	}

	payload := make([]byte, size)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}
	return codec.Wrap(payload, c.order), nil
}

func (c *Client[T]) Send(encoder packet.Encoder) error {
	body := codec.New(c.order)
	body.WriteUInt16(encoder.ID())
	encoder.Encode(body)

	return c.SendRaw(body.Bytes())
}

func (c *Client[T]) SendRaw(data []byte) error {
	frame := codec.New(c.order)
	frame.WriteInt32(int32(len(data)))
	frame.WriteBytes(data)
	_, err := c.conn.Write(frame.Bytes())
	return err
}

func (c *Client[T]) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
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
