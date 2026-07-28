package tcp

import (
	"GoNetworkUtils/codec"
	"GoNetworkUtils/packet"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

type Client[T any] struct {
	data       *T
	server *ServerTcp[T]
	dispatcher *packet.Dispatcher
	byteOrder  binary.ByteOrder
	conn       net.Conn
	packetreceived func(buf *codec.ByteBuf)
}

func CreateClient[T any](data *T, conn net.Conn, order binary.ByteOrder, packetreceived func(buf *codec.ByteBuf)) *Client[T]{
	return &Client[T]{
		data: data,
		byteOrder: order,
		conn: conn,
		packetreceived: packetreceived,
	}
}

func (self *Client[T]) ReadLoop() {
	for {
		buf, err := self.read()
		if err != nil{
			log.Println(err)
			return
		}
		self.packetreceived(buf)
	}
}

func (self *Client[T]) read() (*codec.ByteBuf, error) {
	var header []byte = make([]byte, 4)

	_, err := io.ReadFull(self.conn, header)
	if err != nil {
		return nil, errors.New("Cant Read The Header")
	}
	var headerBuf *codec.ByteBuf = codec.Wrap(header, self.byteOrder)
	var size int32 = headerBuf.ReadInt32()

	var payload []byte = make([]byte, size)
	_, payloadErr := io.ReadFull(self.conn, payload)
	if payloadErr != nil {
		return nil, errors.New("Cant Read The Payload")
	}

	var payloadBuf *codec.ByteBuf = codec.Wrap(payload, self.byteOrder)
	return payloadBuf, nil
}

func (self *Client[T]) Send(packet packet.Packet) {

}

func (self *Client[T]) Close() {
	self.conn.Close()
}
