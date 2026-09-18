package server

import (
	"fmt"
	"httpz/internal/request"
	"httpz/internal/response"
	"net"
	"os"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	open     atomic.Bool
	handler  Handler
}

type Handler func(w *response.Writer, req *request.Request)

func (s *Server) Close() error {
	s.open.Store(false)
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)

	// if hanlderErr != nil {
	// 	err = hanlderErr.WriteError(conn)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	return
	// }

	// b := handlerBody.Bytes()
	// defaultHeders := response.GetDefaultHeaders(len(b))
	// err = response.WriteStatusLine(conn, response.StatusCodeOk)
	// response.WriteHeaders(conn, defaultHeders)

	// _, err = conn.Write(writer.)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
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

// func (e *HandlerError) WriteError(w io.Writer) error {
// 	err := response.WriteStatusLine(w, e.StatusCode)
// 	if err != nil {
// 		return err
// 	}
// 	headers := response.GetDefaultHeaders(len([]byte(e.Message)))
// 	response.WriteHeaders(w, headers)
// 	_, err = w.Write([]byte(e.Message))
// 	return err
// }

func Serve(handler Handler, port int) (*Server, error) {
	fmt.Printf("Starting server on port: %s\n", strconv.Itoa(port))
	listerner, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}
	server := Server{
		handler:  handler,
		listener: listerner,
		open:     atomic.Bool{},
	}
	server.open.Store(true)

	go server.listen()

	return &server, nil
}
