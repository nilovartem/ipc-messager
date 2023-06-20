package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/auth"
)

func TestClient(t *testing.T) {
	c, err := net.Dial("unix", "/tmp/server.sock")
	if err != nil {
		log.Fatal(err)
	}

	//go reader(c)
	certificate := auth.Current()

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(certificate)
	fmt.Printf("Sending := %+v", certificate)
	fmt.Println(buffer.Bytes())
	_, err = c.Write(buffer.Bytes())
	if err != nil {
		log.Fatal("write error:", err)
	}
	c.Close()
	//reader(c)
	//time.Sleep(5 * time.Second)
	//TODO: добавить таймер для таймаута
	//теперь нужно проверить, создался ли сокет /tmp/clientPID.sock . Если создался, то все успешно, если нет - ошибка по timeout
	//if _, err := os.Stat("/tmp/" + certificate.PID.String() + ".sock"); err == nil {
	timer := time.NewTimer(time.Second * 5)
	for {
		select {
		case <-timer.C:
			{
				fmt.Println("Timer fired")
				return
			}
		default:
			{
				time.Sleep(time.Millisecond * 200)
				c, err = net.Dial("unix", "/tmp/"+certificate.PID.String()+".sock")
				if err == nil {
					_, err = c.Write([]byte("Hi!"))
					if err != nil {
						log.Fatal("write error:", err)
					}
					fmt.Println(c.RemoteAddr().String())
					c.Close()
					fmt.Println("Send Hello")
					return
				}
			}
		}
	}
}
func reader(r io.Reader) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		return
	}
	println("Client got:", string(buf[0:n]))
}
