package server

import (
	"bytes"
	"context"
	"net"
	"os"
	"sync"
	"time"

	"github.com/nilovartem/ipc-messager/pkg/auth"
	"github.com/nilovartem/ipc-messager/pkg/client"
)

const (
	DEFAULT_TIMEOUT = time.Duration(time.Millisecond * 100) // 100ms = default timeout for reading and writing
)

// A server is an entity that controls incoming messages and clients
type server struct {
	Certificate auth.Data
	clients     map[auth.PID]client.Client             //
	gateway     string                                 //socket address
	connection  *net.Conn                              //single connection
	handler     func(request []byte) (response []byte) //custom function for handling messages
	close       context.CancelFunc
	errors      chan error    //common channel for errors
	timeout     time.Duration //timeout for reading and writing
	mutex       *sync.RWMutex
}

// Block blocks reading messages from the client
func (s *server) Block(c client.Client) {
	c.Accepted = false
}

// Allow allows reading messages from the client
func (s *server) Allow(c client.Client) {
	c.Accepted = true
}

// New creates new server and returns it
func New(gateway string, timeout time.Duration, handler func(request []byte) (response []byte)) server {
	s := server{
		Certificate: auth.Current(),
		clients:     make(map[auth.PID]client.Client),
		gateway:     gateway,
		handler:     handler,
		connection:  nil,
		close:       nil,
		timeout:     timeout,
		mutex:       &sync.RWMutex{},
	}
	return s
}

// serve reads connection and dispatches clients to handle()
func (s *server) serve(ctx context.Context) {
	select {
	case <-ctx.Done():
		{
			s.cleanup()
			return
		}
	default:
		{
			listener, err := net.Listen("unix", s.gateway)
			if err != nil {
				s.errors <- err
				s.Close()

			}
			defer listener.Close()
			for {
				incoming, err := listener.Accept()
				if err != nil {
					s.errors <- err
				}
				go s.handle(incoming)
			}

		}
	}

}

// lookup checks if the connection with the client is maintained
func (s *server) lookup(pid auth.PID) (ok bool) {
	_, ok = s.clients[pid]
	return
}

// addClient adds client to clients map
func (s *server) addClient(m client.Message) {
	(*s.mutex).Lock()
	s.clients[m.Author.Certificate.PID] = m.Author
	(*s.mutex).Unlock()
}

// Receive reads bytes from connection
func (s *server) Receive(c net.Conn) (bool, []byte) {
	c.SetReadDeadline(time.Now().Add(s.timeout)) //100ms на чтение
	var inc bytes.Buffer
	inc.ReadFrom(c)
	if inc.Len() == 0 {
		return false, inc.Bytes()
	}
	return true, inc.Bytes()
}

// Send sends bytes to connection
func (s *server) Send(c net.Conn, data []byte) {
	c.SetWriteDeadline(time.Now().Add(s.timeout)) //100ms на запись
	inc := bytes.NewBuffer(data)
	inc.WriteTo(c)
}

// handle controls steps of converting the received data to messages and vice versa
func (s *server) handle(c net.Conn) {
	for {
		if ok, b := s.Receive(c); ok {
			m := client.Message{}
			err := m.Unmarshall(b)
			if err != nil {
				s.errors <- err
			}
			if !s.lookup(m.Author.Certificate.PID) {
				s.addClient(m)
			}
			if m.Author.Accepted {
				response := s.handler(m.Content)
				m.CreateMessage(response)
				response, err = m.Marshall()
				if err != nil {
					s.errors <- err
				}
				s.Send(c, response)
			}
		}

	}

}

// Listen starts listening on server and returns channel for all errors that might occur
func (s *server) Listen() <-chan error {
	ctx := context.Background()
	ctx, close := context.WithCancel(ctx)
	s.close = close
	go s.serve(ctx)
	errors := make(chan error)
	s.errors = errors
	return errors
}

// Close initiate closing and cleanup processes
func (s *server) Close() {
	s.close()
	s.cleanup()
}

// cleanup closes server and all connections
func (s *server) cleanup() {
	if s.connection != nil {
		_ = (*s.connection).Close()
	}
	os.RemoveAll(s.gateway)
}
