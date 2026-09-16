package main

import (
	"fmt"
	"httpz/internal/request"
	"maps"
	"net"
	"os"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Accepted connection from", con.RemoteAddr())

		req, err := request.RequestFromReader(con)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for key := range maps.Keys(req.Headers) {
			fmt.Printf("- %s: %s\n", key, req.Headers[key])
		}
		fmt.Printf("Body:\n")
		fmt.Println(string(req.Body))

		fmt.Println("Connection to ", con.RemoteAddr(), "closed")
	}
}
