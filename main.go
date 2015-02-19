package main

import (
	"fmt"
	"github.com/lukad/tracks/server"
	flag "github.com/ogier/pflag"
	"github.com/op/go-logging"
)

// flags
var (
	listen = flag.StringP("listen", "l", ":1337", "Bind to this address")
)

var (
	log    = logging.MustGetLogger("tracks")
	format = logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: tracks [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	logging.SetFormatter(format)
}

func main() {
	s, err := server.Listen(*listen)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	log.Info("Listening on %s\n", s.Addr())

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
