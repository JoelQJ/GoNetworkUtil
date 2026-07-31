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
	opts         *ConnectionOptions[T]
	conn         net.Conn
	byteOrder    binary.ByteOrder
	MaxFrameSize int
}

type ConnectionOptions[T any] struct {
	ByteOrder    binary.ByteOrder
	MaxFrameSize int
	Dispatcher   *packet.Dispatcher
	OnConnect    func(*Client[T])
	OnRawPacket  func(*Client[T], *codec.ByteBuf)
	OnPacket     func(*Client[T], packet.Packet)
	OnDisconnect func(*Client[T])
}

func NewClient[T any](conn net.Conn, data *T) *Client[T] {
	return &Client[T]{
		Data:  data,
		conn:  conn,
	}
}

func Connect[T any](address string, data *T, opts *ConnectionOptions[T]) (*Client[T], error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	client := NewClient(conn, data)
	client.opts = opts
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
	if opts.MaxFrameSize > 0 {
		c.MaxFrameSize = opts.MaxFrameSize
	} else {
		c.MaxFrameSize = DefaultMaxFrameSize
	}
}

func (c *Client[T]) ReadLoop() {
	opts := c.opts
	if opts != nil && opts.OnConnect != nil {
		opts.OnConnect(c)
	}
	defer func() {
		c.Close()
		if opts != nil && opts.OnDisconnect != nil {
			opts.OnDisconnect(c)
		}
	}()

	for {
		buf, err := c.ReadRaw()
		if err != nil {
			return
		}

		if opts != nil && opts.OnRawPacket != nil {
			opts.OnRawPacket(c, buf)
		}

		if opts != nil && opts.Dispatcher != nil {
			if buf.Len() < 2 {
				c.Close()
				return
			}
			packetId := buf.ReadUInt16()
			if pkt, err := opts.Dispatcher.Decode(packetId, buf); err != nil {
				c.Close()
				return
			} else if opts.OnPacket != nil {
				opts.OnPacket(c, pkt)
			}
		}
	}
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

func (c *Client[T]) Send(encoder packet.Encoder) error {
	body := codec.New(c.byteOrder)
	body.WriteUInt16(encoder.Id())
	encoder.Encode(body)

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
