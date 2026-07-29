package tcp

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"github.com/JoelQJ/GoNetworkUtil/packet"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const DefaultMaxFrameSize = 10 * 1024 * 1024

type Client[T any] struct {
	Data         *T
	conn         net.Conn
	byteOrder    binary.ByteOrder
	MaxFrameSize int
}

func NewClient[T any](conn net.Conn, data *T, order binary.ByteOrder) *Client[T] {
	return &Client[T]{
		Data:         data,
		conn:         conn,
		byteOrder:    order,
		MaxFrameSize: DefaultMaxFrameSize,
	}
}

func Connect[T any](address string, data *T, order binary.ByteOrder) (*Client[T], error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, data, order), nil
}

func (c *Client[T]) ReadRaw() (*codec.ByteBuf, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return nil, err
	}

	buf := codec.Wrap(header, c.byteOrder)
	size := buf.ReadInt32()

	if size < 2 || size > int32(c.MaxFrameSize) {
		return nil, errors.New("invalid frame size")
	}

	payload := make([]byte, size)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	return codec.Wrap(payload, c.byteOrder), nil
}

type PacketHandler[T any] interface {
	OnFinishDecode(*Client[T])
}

func (c *Client[T]) Send(pkt packet.Packet) error {
	body := codec.New(c.byteOrder)
	body.WriteUInt16(pkt.Id())
	pkt.Encode(body)

	payload := body.ToBytesSlice()

	frame := codec.New(c.byteOrder)
	frame.WriteInt32(int32(len(payload)))
	frame.WriteBytes(payload)

	_, err := c.conn.Write(frame.ToBytesSlice())
	return err
}

func (c *Client[T]) SendRaw(data []byte) error {
	frame := codec.New(c.byteOrder)
	frame.WriteInt32(int32(len(data)))
	frame.WriteBytes(data)
	_, err := c.conn.Write(frame.ToBytesSlice())
	return err
}

func (c *Client[T]) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client[T]) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Client[T]) Close() error {
	return c.conn.Close()
}
