package server

import (
	"fmt"
	"httpz/internal/response"
	"net"
	"os"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	open     atomic.Bool
}

func (s *Server) Close() error {
	s.open.Store(false)
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	// req, err := request.RequestFromReader(conn)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	defer conn.Close()

	// resp := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n")

	err := response.WriteStatusLine(conn, response.StatusCodeOk)
	defaultHeders := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, defaultHeders); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// _, err = conn.Write(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (s *Server) listen() {
	for {
		if !s.open.Load() {
			fmt.Println("Server is closed.")
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Connection error exiting: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Accepted connection from: %s\n", conn.RemoteAddr())

		go func() {
			s.handle(conn)
		}()
	}
}

func Serve(port int) (*Server, error) {
	fmt.Printf("Starting server on port: %s\n", strconv.Itoa(port))
	listerner, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	server := Server{
		listener: listerner,
		open:     atomic.Bool{},
	}
	server.open.Store(true)

	go server.listen()

	return &server, nil
}
