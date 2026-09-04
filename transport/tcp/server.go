package tcp

import (
	"errors"
	"log"
	"net"
)

func NewServer[T any](factory func(net.Conn) *Client[T]) *Server[T] {
	return &Server[T]{factory: factory}
}

func (s *Server[T]) Bind(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			log.Printf("listen %s on %s: %v", opErr.Op, opErr.Net, opErr.Err)
		}
		return err
	}
	defer listener.Close()
	log.Println("TCP listening on:", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			log.Println(err)
			continue
		}
		go s.factory(conn).ReadLoop()
	}
}
