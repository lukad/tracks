package main

import (
	"fmt"
	"github.com/lukad/tracks/server"
	"os"
)

func main() {
	s, err := server.Listen(":1337")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		os.Exit(1)
	}
	fmt.Printf("Listening on %s\n", s.Addr())

	if err := s.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Scanln()
}
