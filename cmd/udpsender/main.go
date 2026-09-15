package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const remoteAddr = "localhost:42069"

func main() {
	udpAdr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	con, err := net.DialUDP("udp", nil, udpAdr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}

		read, err := con.Write(line)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("UDP Conn wrote: %d\n", read)
	}
}
