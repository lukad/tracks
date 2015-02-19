package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lunixbochs/struc"
	"github.com/op/go-logging"
	"math/rand"
	"net"
	"time"
)

type Server struct {
	conn *net.UDPConn
	log  *logging.Logger
}

// Starts listening on specified address and returns a Server object
func Listen(address string) (*Server, error) {
	addr, _ := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to listen on %s: %s", addr.String(), err))
	}

	server := &Server{
		conn: conn,
		log:  logging.MustGetLogger("server"),
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
			s.log.Debug("Discarding packet smaller than 16 bytes")
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
		s.log.Warning("error unpacking request header: %s", err)
		return
	}

	switch header.Action {

	case actionConnect:
		var req connectRequest
		if err := struc.UnpackWithOrder(buf, &req, binary.BigEndian); err != nil {
			s.log.Warning("Error unpacking connect request: %s", err)
			return
		}
		s.handleConnectRequest(addr, header, req)

	case actionAnnounce:
		var req announceRequest
		if err := struc.UnpackWithOrder(buf, &req, binary.BigEndian); err != nil {
			s.log.Warning("Error unpacking announce request: %s", err)
			return
		}
		s.handleAnnounceRequest(addr, header, req)

	case actionScrape:
		var req scrapeRequest
		if err := struc.UnpackWithOrder(buf, &req, binary.BigEndian); err != nil {
			s.log.Warning("Error unpacking scrape request: %s", err)
			return
		}
		s.handleScrapeRequest(addr, header, req)
	}
}

func (s *Server) handleConnectRequest(addr *net.UDPAddr, header requestHeader, req connectRequest) {
	if header.ConnectionId != 0x41727101980 {
		return
	}
	s.log.Debug("%#v\n", header)
	s.log.Debug("%#v\n", req)

	response := connectResponse{
		Action:        actionConnect,
		TransactionId: req.TransactionId,
		ConnectionId:  rand.Int63(),
	}

	buf := bytes.NewBuffer(nil)
	if err := struc.PackWithOrder(buf, &response, binary.BigEndian); err != nil {
		s.log.Warning("Error packing connect response:", err)
		return
	}

	if _, err := s.conn.WriteToUDP(buf.Bytes(), addr); err != nil {
		s.log.Warning("error:", err)
		return
	}
}

func (s *Server) handleAnnounceRequest(addr *net.UDPAddr, header requestHeader, req announceRequest) {
	s.log.Debug("%#v\n", header)
	s.log.Debug("%#v\n", req)

	peers := []peer{peer{}, peer{}}
	response := announceResponse{
		Action:        actionAnnounce,
		TransactionId: req.TransactionId,
		Interval:      10,
		Leechers:      1337,
		Seeders:       7331,
	}

	buf := bytes.NewBuffer(nil)
	if err := struc.PackWithOrder(buf, &response, binary.BigEndian); err != nil {
		s.log.Warning("Error packing announce response:", err)
		return
	}

	for _, p := range peers {
		if err := struc.PackWithOrder(buf, &p, binary.BigEndian); err != nil {
			s.log.Warning("Error writing peer struct to annnounce response:", err)
			return
		}
	}

	if _, err := s.conn.WriteToUDP(buf.Bytes(), addr); err != nil {
		s.log.Warning("Error sending announce response:", err)
	}
}

func (s *Server) handleScrapeRequest(addr *net.UDPAddr, header requestHeader, req scrapeRequest) {
	s.log.Debug("%#v\n", header)
	s.log.Debug("%#v\n", req)
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
