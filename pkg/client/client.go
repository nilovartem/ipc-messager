package client

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/auth"
)

type Client struct {
	Certificate auth.Data //client's credentials
	Socket      string
	Connection  *net.Conn
	//Accepted    bool //по умолчанию - не обслуживается - нужно выгадать время с таймаутом
}

func (c *Client) accept() { //вернуть канал для чтения сообщений. и для отправки
	//c.accepted = true
	fmt.Printf("Client %+v was accepted", c)
	fmt.Println("Before Dial = ", c.Connection)
	fmt.Println(c.Socket)
	con, err := net.Dial("unix", c.Socket)
	if err == nil {
		fmt.Println(c.Connection)
		//TODO: add error management in net.dial
		c.Connection = &con
	} else {
		fmt.Print(err.Error())
	}

}

// Сообщение, отправляемое клиентом
type message struct {
	certificate auth.Data
	body        []byte
}

func Connect(server string, timeout time.Duration) (Client, error) {
	c := Client{
		Certificate: auth.Current(),
		Socket:      "",
		Connection:  nil,
	}
	//cleanup
	c.auth()
	//TODO: добавить таймер для таймаута
	//теперь нужно проверить, создался ли сокет /tmp/clientPID.sock . Если создался, то все успешно, если нет - ошибка по timeout
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			{
				return c, errors.New("timer fired. cant establish final connection")
			}
		default:
			{
				if _, err := os.Stat("/tmp/" + c.Certificate.PID.String() + ".sock"); err == nil {
					c.Socket = "/tmp/" + c.Certificate.PID.String() + ".sock"
					c.accept()
					return c, nil
				}
			}
		}
	}
}
func (c *Client) auth() {
	conn, err := net.Dial("unix", "/tmp/server.sock")
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(c.Certificate)
	fmt.Printf("Sending := %+v", c.Certificate)
	fmt.Println(buffer.Bytes())
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Fatal("write error:", err)
	}
	conn.Close()
}
func (c *Client) Disconnect() {
	//disconnect
}
