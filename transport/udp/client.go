package udp

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"encoding/binary"
	"net"
)

type Client[T any] struct {
	Data      *T
	conn      *net.UDPConn
	addr      *net.UDPAddr
	byteOrder binary.ByteOrder
}

func NewClient[T any](conn *net.UDPConn, addr *net.UDPAddr, data *T, order binary.ByteOrder) *Client[T] {
	return &Client[T]{
		Data:      data,
		conn:      conn,
		addr:      addr,
		byteOrder: order,
	}
}

func Connect[T any](address string, data *T, order binary.ByteOrder) (*Client[T], error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Client[T]{
		Data:      data,
		conn:      conn,
		addr:      addr,
		byteOrder: order,
	}, nil
}

func (c *Client[T]) ReadRaw() (*codec.ByteBuf, error) {
	buf := make([]byte, 65535)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return codec.Wrap(buf[:n], c.byteOrder), nil
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
	
	_, err := c.conn.Write(packetBuf.ToBytesSlice())
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
	return c.conn.Close()
}
