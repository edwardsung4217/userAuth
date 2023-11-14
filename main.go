package main

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
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop() // Start accepting connections in a goroutine

	// Wait for a signal to gracefully shut down the server
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	// Signal received, shut down the server
	s.Shutdown()

	return nil
}
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		fmt.Println("new connection to the server:", conn.RemoteAddr())
		go s.readLoop(conn)
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
	close(s.quitch) // Close the quitch channel to unblock Start() method
	s.ln.Close()    // Close the listener to stop accepting new connections
	// You can add additional cleanup code here if necessary
}
func main() {
	server := NewServer(":3000")
	go func() {
		for payload := range server.msgch {
			fmt.Printf("Received from %s: %s\n", payload.user, string(payload.message))
		}
	}()
	if err := server.Start(); err != nil {
		fmt.Println("Server error:", err)
	}
}
