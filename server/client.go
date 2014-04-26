package server

import (
	"net"
)

type Client struct {
	Id   int64
	Addr *net.UDPAddr
}
