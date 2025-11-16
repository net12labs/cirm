package socketserver

import (
	"fmt"
	"net"
	"os"

	"github.com/net12labs/cirm/dali/rtm"
)

type SocketServer struct {
	socketPath string
	listener   net.Listener
	OnMessage  func(data []byte, conn net.Conn)
}

func NewSocketServer() *SocketServer {
	return &SocketServer{
		OnMessage: defaultMessageHandler,
	}
}

func defaultMessageHandler(data []byte, conn net.Conn) {
	fmt.Printf("Received on socket: %s\n", string(data))
	response := "ACK"
	if _, err := conn.Write([]byte(response)); err != nil {
		fmt.Println("Socket write error:", err)
	}
}

func (s *SocketServer) Start(socketPath string) error {
	s.socketPath = socketPath
	rtm.Do.InitFsPath(socketPath)

	// Remove existing socket file if it exists
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing socket: %w", err)
	}

	// Create Unix domain socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}

	s.listener = listener

	// Cleanup socket on exit
	rtm.Runtime.OnExit.AddListener(func(code any) {
		s.Stop()
	})

	// Start listening for connections in a goroutine
	go func() {
		fmt.Println("Socket listening on:", socketPath)
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Check if listener was closed intentionally
				if s.listener == nil {
					return
				}
				return
			}

			// Handle connection in a separate goroutine
			go s.handleConnection(conn)
		}
	}()

	return nil
}

func (s *SocketServer) Stop() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
	if s.socketPath != "" {
		os.Remove(s.socketPath)
	}
}

func (s *SocketServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Socket read error:", err)
		return
	}

	// Process the received data
	data := buf[:n]
	if s.OnMessage != nil {
		s.OnMessage(data, conn)
	}
}
