package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/client"
)

func main() {
	c, err := client.Connect("/tmp/server.sock", time.Millisecond*100)
	fmt.Println("Client", c)
	if err == nil {
		i := 0
		for i < 10 {
			if c.Connection == nil {
				//fmt.Println("AAA NIL CONNECTION")
			} else {
				//var buffer bytes.Buffer = *bytes.NewBufferString("Hi" + strconv.Itoa(i))
				conn := *c.Connection
				//buffer.WriteTo(*c.Connection)
				_, err = conn.Write([]byte("Hi!" + strconv.Itoa(i)))
				(*c.Connection).SetWriteDeadline(time.Now())
				//conn.Close()
				//*c.Connection.
				i++
			}
		}
	}
}
