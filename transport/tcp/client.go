package tcp

import (
	"GoNetworkUtils/codec"
	"GoNetworkUtils/packet"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Client[T any] struct {
	Data       *T
	conn       net.Conn
	dispatcher *packet.Dispatcher
	byteOrder  binary.ByteOrder
}

func NewClient[T any](conn net.Conn, data *T, dispatcher *packet.Dispatcher, order binary.ByteOrder) *Client[T] {
	return &Client[T]{
		Data:       data,
		conn:       conn,
		dispatcher: dispatcher,
		byteOrder:  order,
	}
}

func Connect[T any](address string, data *T, dispatcher *packet.Dispatcher, order binary.ByteOrder) (*Client[T], error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, data, dispatcher, order), nil
}

func (c *Client[T]) ReadLoop(onPacket func(*Client[T], *packet.Packet)) {
	for {
		pkt, err := c.readPacket()
		if err != nil {
			log.Println(err)
			return
		}
		pkt.OnFinishDecode()
		onPacket(c, &pkt)
	}
}

func (c *Client[T]) readPacket() (packet.Packet, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return nil, err
	}

	headerBuf := codec.Wrap(header, c.byteOrder)
	size := headerBuf.ReadInt32()

	payload := make([]byte, size)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	payloadBuf := codec.Wrap(payload, c.byteOrder)
	pkt, err := c.dispatcher.Decode(payloadBuf)
	if err != nil {
		return nil, err
	}
	return pkt, nil
}

func (c *Client[T]) Send(pkt packet.Packet) error {
	buf := codec.New(c.byteOrder)
	buf.WriteUInt16(pkt.Id())
	pkt.Encode(buf)

	payload := buf.ToBytesSlice()

	frame := codec.New(c.byteOrder)
	frame.WriteInt32(int32(len(payload)))
	frame.WriteBytes(payload)

	_, err := c.conn.Write(frame.ToBytesSlice())
	return err
}

func (c *Client[T]) SendRaw(data []byte) error {
	_, err := c.conn.Write(data)
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
