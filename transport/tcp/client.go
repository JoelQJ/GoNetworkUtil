package tcp

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"encoding/binary"
	"io"
	"net"
)

type Client[T any] struct {
	Data      *T
	conn      net.Conn
	byteOrder binary.ByteOrder
}

func NewClient[T any](conn net.Conn, data *T, order binary.ByteOrder) *Client[T] {
	return &Client[T]{
		Data:      data,
		conn:      conn,
		byteOrder: order,
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

	payload := make([]byte, size)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	return codec.Wrap(payload, c.byteOrder), nil
}

type Packet interface {
	Id() uint16
}

type Encoder interface {
	Encode(*codec.ByteBuf) error
}

func (c *Client[T]) Send(pkt Packet, enc Encoder) error {
	var packetBuf *codec.ByteBuf = codec.New(c.byteOrder)
	packetBuf.WriteUInt16(pkt.Id())
	enc.Encode(packetBuf)
	var payload []byte = packetBuf.ToBytesSlice()
	
	var frameBuf *codec.ByteBuf = codec.New(c.byteOrder)
	frameBuf.WriteInt32(int32(len(payload)))
	frameBuf.WriteBytes(payload)
	
	_, err := c.conn.Write(frameBuf.ToBytesSlice())
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
