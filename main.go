package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 8)
	for {
		_, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}

			os.Exit(1)
		}

		fmt.Printf("read: %s\n", data)
	}
}
