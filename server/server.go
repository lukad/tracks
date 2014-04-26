package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Server struct {
	conn        *net.UDPConn
	clients     map[int64]*Client
	clientMutex sync.Mutex
}

// Starts listening on specified address and returns a Server object
func Listen(address string) (*Server, error) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New("Failed to listen on " + addr.String())
	}

	server := &Server{
		conn:    conn,
		clients: make(map[int64]*Client),
	}

	return server, nil
}

func (s *Server) Run() (err error) {
	for {
		b := make([]byte, 1024)
		var n int
		var addr *net.UDPAddr

		n, addr, _ = s.conn.ReadFromUDP(b)
		// TODO: handle error

		if n < 16 { // If packet is smaller than 16 bytes we may ignore it
			continue
		}
		go s.handleRequest(addr, n, b)
	}
	return
}

func (s *Server) handleRequest(addr *net.UDPAddr, n int, b []byte) {
	var header requestHeader
	buf := bytes.NewReader(b)

	if err := binary.Read(buf, binary.BigEndian, &header); err != nil {
		return
	}

	switch header.Action {

	case actionConnect:
		var req connectRequest
		if err := binary.Read(buf, binary.BigEndian, &req); err != nil {
			return
		}
		s.handleConnectRequest(addr, header, req)

	case actionAnnounce:
		var req announceRequest
		if err := binary.Read(buf, binary.BigEndian, &req); err != nil {
			return
		}
		s.handleAnnounceRequest(addr, header, req)
	}

}

func (s *Server) handleConnectRequest(addr *net.UDPAddr, header requestHeader, req connectRequest) {
	if header.ConnectionId != 0x41727101980 {
		return
	}
	response := connectResponse{
		Action:        actionConnect,
		TransactionId: req.TransactionId,
		ConnectionId:  header.ConnectionId,
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, response)
	if n, err := s.conn.WriteToUDP(buf.Bytes(), addr); err != nil || n != 16 {
		return
	}
	c := &Client{
		Addr: addr,
		Id:   rand.Int63(),
	}
	s.addClient(c)
}

func (s *Server) handleAnnounceRequest(addr *net.UDPAddr, header requestHeader, req announceRequest) {

}

func (s *Server) addClient(c *Client) {
	s.clientMutex.Lock()
	s.clients[c.Id] = c
	s.clientMutex.Unlock()
}

func (s *Server) removeClient(id int64) {
	s.clientMutex.Lock()
	delete(s.clients, id)
	s.clientMutex.Unlock()
}

func (s *Server) Addr() *net.UDPAddr {
	return s.conn.LocalAddr().(*net.UDPAddr)
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
