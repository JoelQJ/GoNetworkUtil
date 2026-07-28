package udp

import (
	"GoNetworkUtils/codec"
	"GoNetworkUtils/packet"
	"encoding/binary"
	"log"
	"net"
	"sync"
)

type ServerUdp[T any] struct {
	conn          *net.UDPConn
	clients       map[string]*Client[T]
	clientFactory func(addr *net.UDPAddr) *Client[T]
	byteOrder     binary.ByteOrder
	OnConnect     func(*Client[T])
	OnPacket      func(*Client[T], packet.Packet)
	OnDisconnect  func(*Client[T])
	mu            sync.RWMutex
}

func NewServer[T any](clientFactory func(addr *net.UDPAddr) *Client[T], order binary.ByteOrder) *ServerUdp[T] {
	return &ServerUdp[T]{
		clients:       make(map[string]*Client[T]),
		clientFactory: clientFactory,
		byteOrder:     order,
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

		addrStr := remoteAddr.String()
		client := s.getOrCreateClient(addrStr, remoteAddr)

		payloadBuf := codec.Wrap(buf[:n], s.byteOrder)
		pkt, err := client.dispatcher.Decode(payloadBuf)
		if err != nil {
			log.Println(err)
			continue
		}
		pkt.OnFinishDecode()

		if s.OnPacket != nil {
			s.OnPacket(client, pkt)
		}
	}
}

func (s *ServerUdp[T]) getOrCreateClient(addrStr string, remoteAddr *net.UDPAddr) *Client[T] {
	s.mu.RLock()
	client, exists := s.clients[addrStr]
	s.mu.RUnlock()
	if exists {
		return client
	}

	client = s.clientFactory(remoteAddr)
	s.mu.Lock()
	s.clients[addrStr] = client
	s.mu.Unlock()

	if s.OnConnect != nil {
		s.OnConnect(client)
	}
	return client
}

func (s *ServerUdp[T]) RemoveClient(client *Client[T]) {
	s.mu.Lock()
	delete(s.clients, client.addr.String())
	s.mu.Unlock()

	if s.OnDisconnect != nil {
		s.OnDisconnect(client)
	}
}

func (s *ServerUdp[T]) Broadcast(pkt packet.Packet) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, client := range s.clients {
		client.Send(pkt)
	}
}

func (s *ServerUdp[T]) ClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

func (s *ServerUdp[T]) GetClients() []*Client[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	clients := make([]*Client[T], 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	return clients
}
