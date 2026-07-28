package tcp

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"errors"
	"log"
	"net"
	"strconv"
)

type ServerTcp[T any] struct {
	clientFactory func(conn net.Conn) *Client[T]
	OnConnect     func(*Client[T])
	OnRawPacket   func(*Client[T], *codec.ByteBuf)
	OnDisconnect  func(*Client[T])
}

func NewServer[T any](clientFactory func(conn net.Conn) *Client[T]) *ServerTcp[T] {
	return &ServerTcp[T]{
		clientFactory: clientFactory,
	}
}

func (s *ServerTcp[T]) Bind(ipInterface string, port int64) {
	address := ipInterface + ":" + strconv.FormatInt(port, 10)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		if opErr, ok := errors.AsType[*net.OpError](err); ok {
			log.Println("Operación:", opErr.Op)
			log.Println("Red:", opErr.Net)
			log.Println("Error:", opErr.Err)
		}
		return
	}
	defer listener.Close()
	log.Println("Server Listening on:", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		client := s.clientFactory(conn)
		if s.OnConnect != nil {
			s.OnConnect(client)
		}
		go s.handleClient(client)
	}
}

func (s *ServerTcp[T]) handleClient(client *Client[T]) {
	defer func() {
		client.Close()
		if s.OnDisconnect != nil {
			s.OnDisconnect(client)
		}
	}()

	for {
		buf, err := client.ReadRaw()
		if err != nil {
			return
		}

		if s.OnRawPacket != nil {
			s.OnRawPacket(client, buf)
		}
	}
}
