package server

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	s := New("/tmp/server.sock")
	err := s.Listen()
	if err == nil {
		for client := range s.line {
			fmt.Println("TestServer = connected client ", client.certificate)
			client.Accept()
		}
	}
}
