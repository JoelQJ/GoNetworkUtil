package udp

import (
	"log"
	"net"

	"github.com/JoelQJ/GoNetworkUtil/codec"
)

func NewServer[T any](factory func(*net.UDPAddr) *Client[T], opts *ConnectionOptions[T]) *Server[T] {
	return &Server[T]{factory: factory, opts: opts}
}

func (s *Server[T]) Bind(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.conn = conn
	defer conn.Close()
	log.Println("UDP listening on:", addr)

	buf := make([]byte, 65535)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		client := s.factory(remoteAddr)
		if s.opts != nil && s.opts.OnConnect != nil {
			s.opts.OnConnect(client)
		}

		payload := codec.Wrap(buf[:n], client.order)
		if s.opts != nil && s.opts.OnRawPacket != nil {
			s.opts.OnRawPacket(client, payload)
		}
	}
}
