package udp

import (
	"GoNetworkUtils/codec"
	"GoNetworkUtils/packet"
	"encoding/binary"
	"log"
	"net"
)

type Client[T any] struct {
	Data       *T
	conn       *net.UDPConn
	addr       *net.UDPAddr
	dispatcher *packet.Dispatcher
	byteOrder  binary.ByteOrder
}

func NewClient[T any](conn *net.UDPConn, addr *net.UDPAddr, data *T, dispatcher *packet.Dispatcher, order binary.ByteOrder) *Client[T] {
	return &Client[T]{
		Data:       data,
		conn:       conn,
		addr:       addr,
		dispatcher: dispatcher,
		byteOrder:  order,
	}
}

func Connect[T any](address string, data *T, dispatcher *packet.Dispatcher, order binary.ByteOrder) (*Client[T], error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Client[T]{
		Data:       data,
		conn:       conn,
		addr:       addr,
		dispatcher: dispatcher,
		byteOrder:  order,
	}, nil
}

func (c *Client[T]) ReadLoop(onPacket func(*Client[T], packet.Packet)) {
	buf := make([]byte, 65535)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		payloadBuf := codec.Wrap(buf[:n], c.byteOrder)
		pkt, err := c.dispatcher.Decode(payloadBuf)
		if err != nil {
			log.Println(err)
			continue
		}
		pkt.OnFinishDecode()
		onPacket(c, pkt)
	}
}

func (c *Client[T]) Send(pkt packet.Packet) error {
	buf := codec.New(c.byteOrder)
	buf.WriteUInt16(pkt.Id())
	pkt.Encode(buf)
	_, err := c.conn.Write(buf.ToBytesSlice())
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
