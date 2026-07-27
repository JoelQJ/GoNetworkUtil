package tcp

import (
	"GoNetworkUtils/packet"
	"net"
)

type Client[T any] struct{
	data *T
	conn net.Conn
}

func (c *Client[T]) ReadLoop(){
	for{
		
	}
}

func (c *Client[T]) Send(packet packet.Packet){

}

func (c *Client[T]) Close(){
	c.conn.Close()
}