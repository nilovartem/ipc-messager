package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/nilovartem/ipc-messager/pkg/auth"
	"github.com/nilovartem/ipc-messager/pkg/client"
)

type server struct {
	Certificate auth.Data
	Clients     map[auth.PID]client.Client
	gateway     string              //special socket for getting clients in
	line        *chan client.Client //очередь клиентов, которые хотят подключиться к серверу. Их надо либо принять, либо откинуть
	close       chan struct{}
}

func (s *server) Accept(c client.Client) { //вернуть канал для чтения сообщений. и для отправки
	//c.Accepted = true
	fmt.Printf("Client %+v was accepted", c)
	listener, _ := net.Listen("unix", c.Socket)
	*c.Connection, _ = listener.Accept()

}

func New(gateway string) server {
	fmt.Println("New server")
	s := server{
		Certificate: auth.Current(),
		Clients:     make(map[auth.PID]client.Client),
		gateway:     gateway,
		line:        nil, //небуферизованный канал
		close:       make(chan struct{}, 1),
	}
	s.cleanup() //TODO: может и не надо чистить
	return s
}

// Проверяет входящих клиентов на наличие сертификата и ставит в очередь
func (s *server) serveGateway() {
	listener, err := net.Listen("unix", s.gateway)
	//fmt.Println(listener.Addr().String())
	if err != nil {
		//return err
	}
	defer listener.Close()
	for {
		//<-s.close
		// Accept new connections, dispatching them to echoServer
		// in a goroutine.
		incoming, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		//первое сообщение, отправленное клиентом - сертификат
		if valid, c, err := verifyClient(incoming); valid {
			fmt.Println("valid")
			//добавляем клиента в список активных
			s.Clients[c.Certificate.PID] = c
			fmt.Println("gateway client address ", &c)
			*s.line <- c
			fmt.Println("Line length ", len(*s.line))
		} else {
			log.Fatal("client error:", err)
		}
		incoming.Close()
		/*client := Client{
			certificate: ,
		}*/

		go echoServer(incoming)
	}
}

// канал для прослушивания и ошибка и может быть канал для ошибок
func (s *server) Listen() (<-chan client.Client, error) {
	//fmt.Println("listen")
	//line := make(chan client) //очередь клиентов, которые хотят подключиться к серверу. Их надо либо принять, либо откинуть
	go s.serveGateway()
	line := make(chan client.Client)
	s.line = &line
	return line, nil
}
func (s *server) Close() {
	s.close <- struct{}{}
	close(s.close)
}
func (s *server) cleanup() {
	if err := os.RemoveAll(s.gateway); err != nil {
		log.Fatal(err)
	}
}

// проверить клиента на валидность сертификата
func verifyClient(incoming net.Conn) (bool, client.Client, error) {

	//Шаблон для клиента
	c := client.Client{
		Certificate: auth.Data{},
		Socket:      "",
		Connection:  nil,
		//messages:    make(chan interface{}),
	}

	var buffer bytes.Buffer
	buffer.ReadFrom(incoming)
	certificate := auth.Data{}
	enc := gob.NewDecoder(&buffer)
	err := enc.Decode(&certificate)

	//fmt.Printf("Incoming := %+v", certificate)
	//fmt.Println(buffer.Bytes())

	if err != nil {
		return false, c, err
	}
	c.Certificate = certificate
	c.Connection = &incoming
	//TODO: make good path and make it simple
	c.Socket = "/tmp/" + c.Certificate.PID.String() + ".sock"
	//fmt.Printf("All good : %+v", c)
	return true, c, err
}
func echoServer(c net.Conn) {
	log.Printf("Client connected [%s]", c.RemoteAddr().Network())
	io.Copy(c, c)
	c.Close()
}
