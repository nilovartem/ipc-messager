package client

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/auth"
)

const (
	DEFAULT_TIMEOUT = time.Duration(time.Millisecond * 100) // 100ms = default timeout for reading and writing
)

// A message contains content and author of received or sended message
type Message struct {
	Author  Client
	Content []byte
}

// Method CreateMessage creates new Message and implementing IMessage
func (m *Message) CreateMessage(data []byte) {
	c := Client{
		Certificate: auth.Current(), //подставляем данные сервера
		Accepted:    true,
		connection:  nil,
	}
	m.Author = c
	m.Content = data
}

// Method Marshall encodes Message to byte array
func (m *Message) Marshall() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Method Unmarshall decodes byte array to Message
func (m *Message) Unmarshall(request []byte) error {
	c := Client{
		Certificate: auth.Data{},
		Accepted:    true,
	}
	m.Author = c
	var buffer bytes.Buffer = *bytes.NewBuffer(request)
	enc := gob.NewDecoder(&buffer)
	err := enc.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

// A Client represents a connection to server
type Client struct {
	Certificate auth.Data //client's credentials
	Accepted    bool      //по умолчанию - не обслуживается - нужно выгадать время с таймаутом
	connection  net.Conn
	timeout     time.Duration
}

// Method Send sends byte array to connection in Client
func (c *Client) Send(data []byte) error {
	c.connection.SetWriteDeadline(time.Now().Add(c.timeout)) //100ms на запись
	m := Message{}
	m.CreateMessage(data)
	data, err := m.Marshall()
	if err != nil {
		return err
	}
	inc := bytes.NewBuffer(data)
	inc.WriteTo(c.connection)
	return nil
}

// Method Receive receives byte array from connection in Client
func (c *Client) Receive() ([]byte, bool) {
	c.connection.SetReadDeadline(time.Now().Add(c.timeout)) //100ms на чтение
	var b []byte
	var inc bytes.Buffer
	inc.ReadFrom(c.connection)
	b = inc.Bytes()
	if inc.Len() == 0 {
		return nil, false
	}
	m := Message{}
	err := m.Unmarshall(b)
	if err != nil {
		fmt.Println(err)
	}
	return m.Content, true

}

// Connect establish connection to server within timeout
func Connect(server string, timeout time.Duration) (Client, error) {
	c := Client{
		Certificate: auth.Current(),
		Accepted:    true,
		timeout:     timeout,
	}
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			{
				return c, errors.New("can't connect to server")
			}
		default:
			{
				conn, _ := net.Dial("unix", server)
				if conn != nil {
					c.connection = conn
					return c, nil
				}
			}
		}
	}
}

// Disconnect closes connection with server
func (c *Client) Disconnect() {
	if c.connection != nil {
		c.connection.Close()
	}
}
