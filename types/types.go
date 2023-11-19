package types

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Payload struct {
	user    string
	message []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Payload
	shutdown   bool
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Payload, 10),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.acceptLoop() // Start accepting connections in a goroutine

	// Wait for a signal to gracefully shut down the server
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	// Signal received, shut down the server
	s.Shutdown()

	// Close the listener after the server has shut down
	defer ln.Close()

	return nil
}

func (s *Server) acceptLoop() {
	for {
		select {
		case <-s.quitch:
			return
		default:
			conn, err := s.ln.Accept()
			if err != nil {
				if s.shutdown {
					return
				}

				if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
					continue
				}
				fmt.Println("Accept error:", err)
				return
			}
			fmt.Println("new connection to the server:", conn.RemoteAddr())
			go s.readLoop(conn)
		}
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
				return
			}
			fmt.Println("read error:", err)
			return
		}
		// msg := buf[:n]
		// fmt.Println(string(msg))
		s.msgch <- Payload{
			user:    conn.RemoteAddr().String(),
			message: buf[:n],
		}
	}
}
func (s *Server) Shutdown() {
	if s.shutdown {
		return // Server has already been shut down
	}
	s.shutdown = true

	close(s.quitch) // Close the quitch channel to unblock Start() method

	// Close the listener only if it is not nil
	if s.ln != nil {
		s.ln.Close()
	}
	// You can add additional cleanup code here if necessary
}

func (s *Server) GetMsgChannel() <-chan Payload {
	return s.msgch
}
func (p *Payload) GetUser() string {
	return p.user
}
func (p *Payload) GetMessage() []byte {
	return p.message
}
