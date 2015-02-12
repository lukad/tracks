package main

import (
	"fmt"
	"github.com/lukad/tracks/server"
	flag "github.com/ogier/pflag"
	"log"
)

// flags
var (
	listen = flag.StringP("listen", "l", ":1337", "Bind to this address")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: tracks [options]")
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetPrefix("[tracks] ")
}

func main() {
	s, err := server.Listen(*listen)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	log.Printf("Listening on %s\n", s.Addr())

	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
