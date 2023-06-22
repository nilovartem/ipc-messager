package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/server"
)

// custom function for message handling
func handler(request []byte) (response []byte) {
	var buffer bytes.Buffer = *bytes.NewBuffer(request)
	fmt.Println("Это запрос: ", buffer.String())
	value, _ := strconv.Atoi(buffer.String())
	value *= 2
	return []byte(strconv.Itoa(value))

}

// main shows the next example: server running for 10s, than closes
func main() {
	s := server.New("/tmp/server.sock", server.DEFAULT_TIMEOUT, handler)
	timer := time.NewTimer(time.Second * 10)
	select {
	case <-timer.C:
		{
			fmt.Println("Таймер истек")
			s.Close()
		}
	case err := <-s.Listen():
		{
			fmt.Println("Возникла ошибка:", err)
		}
	}
}
