package main

import (
	"bytes"
	"fmt"

	"github.com/nilovartem/ipc-messager/pkg/client"
	"github.com/nilovartem/ipc-messager/pkg/server"
)

func main() {
	s := server.New("/tmp/server.sock")
	line, err := s.Listen()
	if err == nil {
		for c := range line {
			//чекаю клиента, если понравился - принимаю
			s.Accept(c)
			go func (){
				answer := c.Read()
				if answer == lalal
				c.Send(;,g;,;dlfg;lfd)
			}()
			/*
			go func(cl client.Client) {
				fmt.Println("go func server.client")
				for {
					var buffer bytes.Buffer
					buffer.ReadFrom(*cl.Connection)
					if buffer.Len() != 0 {
						fmt.Println(buffer.String())
						fmt.Println(buffer.Len())
					}

				}
			}(c)
*/
		}
	}
}
