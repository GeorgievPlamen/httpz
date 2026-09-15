package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

		lines := getLinesChannel(con)
		for v := range lines {
			fmt.Printf("%s\n", v)
		}
		fmt.Println("Connection to ", con.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)

	go func() {
		defer f.Close()
		defer close(linesChan)

		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					linesChan <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				fmt.Printf("read: %s%s\n", currentLineContents, parts[i])
				linesChan <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return linesChan
}
