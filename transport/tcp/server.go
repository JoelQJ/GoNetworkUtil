package tcp

import (
	"errors"
	"log"
	"net"
	"strconv"
)

type ServerTcp[T any] struct {
	clientFactory func(conn net.Conn) *Client[T]
	opts          *ConnectionOptions[T]
}

func NewServer[T any](clientFactory func(conn net.Conn) *Client[T], opts *ConnectionOptions[T]) *ServerTcp[T] {
	return &ServerTcp[T]{
		clientFactory: clientFactory,
		opts:          opts,
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
		ApplyOptions(client, s.opts)
		if s.opts != nil && s.opts.OnConnect != nil {
			s.opts.OnConnect(client)
		}
		go s.handleClient(client)
	}
}

func (s *ServerTcp[T]) handleClient(client *Client[T]) {
	opts := s.opts
	defer func() {
		client.Close()
		if opts != nil && opts.OnDisconnect != nil {
			opts.OnDisconnect(client)
		}
	}()

	for {
		buf, err := client.ReadRaw()
		if err != nil {
			return
		}

		if opts != nil && opts.OnRawPacket != nil {
			opts.OnRawPacket(client, buf)
		}
	}
}
