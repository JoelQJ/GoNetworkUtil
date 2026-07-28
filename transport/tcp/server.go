package tcp

import (
	"GoNetworkUtils/packet"
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
)

type ServerTcp[T any] struct {
	validConnections map[*Client[T]]struct{}
	clientFactory    func(conn net.Conn) *Client[T]
	OnConnect        func(*Client[T])
	OnPacket         func(*Client[T], *packet.Packet)
	OnDisconnect     func(*Client[T])
	mu               sync.RWMutex
}

func NewServer[T any](clientFactory func(conn net.Conn) *Client[T]) *ServerTcp[T] {
	return &ServerTcp[T]{
		validConnections: make(map[*Client[T]]struct{}),
		clientFactory:    clientFactory,
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
		s.AddClient(client)
		if s.OnConnect != nil {
			s.OnConnect(client)
		}
		go s.handleClient(client)
	}
}

func (s *ServerTcp[T]) handleClient(client *Client[T]) {
	defer func() {
		s.RemoveClient(client)
		client.Close()
		if s.OnDisconnect != nil {
			s.OnDisconnect(client)
		}
	}()

	client.ReadLoop(func(c *Client[T], pkt *packet.Packet) {
		if s.OnPacket != nil {
			s.OnPacket(c, pkt)
		}
	})
}

func (s *ServerTcp[T]) AddClient(client *Client[T]) {
	s.mu.Lock()
	s.validConnections[client] = struct{}{}
	s.mu.Unlock()
}

func (s *ServerTcp[T]) RemoveClient(client *Client[T]) {
	s.mu.Lock()
	delete(s.validConnections, client)
	s.mu.Unlock()
}

func (s *ServerTcp[T]) Broadcast(pkt packet.Packet) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for client := range s.validConnections {
		client.Send(pkt)
	}
}

func (s *ServerTcp[T]) ClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.validConnections)
}

func (s *ServerTcp[T]) GetClients() []*Client[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	clients := make([]*Client[T], 0, len(s.validConnections))
	for client := range s.validConnections {
		clients = append(clients, client)
	}
	return clients
}
