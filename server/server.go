package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/lunixbochs/struc"
	"log"
	"math/rand"
	"net"
	"time"
)

type Server struct {
	conn *net.UDPConn
}

// Starts listening on specified address and returns a Server object
func Listen(address string) (*Server, error) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New("Failed to listen on " + addr.String())
	}

	server := &Server{
		conn: conn,
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
			log.Println("packet is smaller than 16 bytes")
			continue
		}
		go s.handleRequest(addr, n, b)
	}
	return
}

func (s *Server) handleRequest(addr *net.UDPAddr, n int, b []byte) {
	var header requestHeader
	buf := bytes.NewReader(b)

	if err := struc.UnpackWithOrder(buf, &header, binary.BigEndian); err != nil {
		log.Println("error unpacking packet:", err)
		return
	}

	switch header.Action {

	case actionConnect:
		var req connectRequest
		if err := struc.UnpackWithOrder(buf, &req, binary.BigEndian); err != nil {
			log.Println("Error unpacking connect request:", err)
			return
		}
		s.handleConnectRequest(addr, header, req)

	case actionAnnounce:
		var req announceRequest
		if err := struc.UnpackWithOrder(buf, &req, binary.BigEndian); err != nil {
			log.Println("Error unpacking announce request:", err)
			return
		}
		s.handleAnnounceRequest(addr, header, req)
	}
}

func (s *Server) handleConnectRequest(addr *net.UDPAddr, header requestHeader, req connectRequest) {
	if header.ConnectionId != 0x41727101980 {
		return
	}
	log.Printf("%#v\n", header)
	log.Printf("%#v\n", req)

	response := connectResponse{
		Action:        actionConnect,
		TransactionId: req.TransactionId,
		ConnectionId:  rand.Int63(),
	}

	buf := bytes.NewBuffer(nil)
	if err := struc.PackWithOrder(buf, &response, binary.BigEndian); err != nil {
		log.Println("Error packing connect response struct:", err)
		return
	}

	if n, err := s.conn.WriteToUDP(buf.Bytes(), addr); err != nil || n != len(buf.Bytes()) {
		log.Println("error:", err)
		return
	}
}

func (s *Server) handleAnnounceRequest(addr *net.UDPAddr, header requestHeader, req announceRequest) {
	log.Printf("%#v\n", header)
	log.Printf("%#v\n", req)
}

func (s *Server) Addr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
