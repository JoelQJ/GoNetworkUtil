package udp

import (
	"github.com/JoelQJ/GoNetworkUtil/codec"
	"log"
	"net"
)

type ServerUdp[T any] struct {
	conn          *net.UDPConn
	clientFactory func(addr *net.UDPAddr) *Client[T]
	opts          *ConnectionOptions[T]
}

func NewServer[T any](clientFactory func(addr *net.UDPAddr) *Client[T], opts *ConnectionOptions[T]) *ServerUdp[T] {
	return &ServerUdp[T]{
		clientFactory: clientFactory,
		opts:          opts,
	}
}

func (s *ServerUdp[T]) Bind(ipInterface string, port int64) {
	addr := &net.UDPAddr{IP: net.ParseIP(ipInterface), Port: int(port)}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	s.conn = conn
	log.Println("Server UDP Listening on:", addr.String())

	buf := make([]byte, 65535)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		client := s.clientFactory(remoteAddr)
		ApplyOptions(client, s.opts)
		if s.opts != nil && s.opts.OnConnect != nil {
			s.opts.OnConnect(client)
		}

		payload := codec.Wrap(buf[:n], client.byteOrder)

		if s.opts != nil && s.opts.OnRawPacket != nil {
			s.opts.OnRawPacket(client, payload)
		}
	}
}
