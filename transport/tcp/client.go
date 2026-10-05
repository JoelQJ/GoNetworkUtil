package tcp

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"

	"github.com/JoelQJ/GoNetworkUtil/codec"
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
		if c.opts.Dispatcher != nil {
			if err := c.opts.Dispatcher.Dispatch(c, buf); err != nil {
				log.Println(err)
			}
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

func (c *Client[T]) Send(id uint16, buf *codec.ByteBuf) error {
	payload := buf.Bytes()

	frame := codec.New(c.order)
	frame.WriteInt32(int32(2 + len(payload)))
	frame.WriteUInt16(id)
	frame.WriteBytes(payload)

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
	c.CloseWithError(nil)
	return c.closeErr
}

func (c *Client[T]) Err() error {
	return c.closeErr
}

func (c *Client[T]) CloseWithError(err error) {
	c.once.Do(func(){
		c.conn.Close()
		if c.closeErr == nil {
			c.closeErr = err
		}
	})
	

}
