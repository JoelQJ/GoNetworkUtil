package tcp

import (
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
)


type ServerTcp[T any] struct{
	validConnections map[*Client[T]]struct{}
	clientFactory func(conn net.Conn) *Client[T]
	mu sync.Mutex
}

func NewServer[T any](clientFactory func(conn net.Conn) *Client[T]) *ServerTcp[T]{
	return &ServerTcp[T]{clientFactory: clientFactory}
}


func (self *ServerTcp[T]) Bind(ipInterface string, port int64){
	var address string = ipInterface + ":" + strconv.FormatInt(port, 10)
	server, err := net.Listen("tcp", address)
	if err != nil {
		if opErr, ok := errors.AsType[*net.OpError](err); ok{
			log.Println("Operación:", opErr.Op)
        	log.Println("Red:", opErr.Net)
         	log.Println("Error:", opErr.Err)
		}
		return
	}

	defer server.Close()
	log.Println("Server Listening on: ", address)

	for{
		conn, err := server.Accept()
		if err != nil{
			log.Println(err)
			continue
		}
		var client *Client[T] = self.clientFactory(conn)
		go client.ReadLoop()
	}
}

func (self *ServerTcp[T]) AddClient(client *Client[T]){
	self.mu.Lock()
	self.validConnections[client] = struct{}{}
	self.mu.Unlock()
}

func (self *ServerTcp[T]) RemoveClient(client *Client[T]){
	self.mu.Lock()
	self.validConnections[client] = struct{}{}
	self.mu.Unlock()
}
