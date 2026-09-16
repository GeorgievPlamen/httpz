package headers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const Clrf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endIndex := bytes.Index(data, []byte(Clrf))
	if endIndex == -1 {
		return n, done, nil
	}

	if endIndex == 0 {
		return len([]byte(Clrf)), true, nil
	}

	buff := make([]byte, 1024)
	reader := bytes.NewReader(data[:endIndex])
	read, err := reader.Read(buff)
	if err != nil {
		if errors.Is(err, io.EOF) {
			done = true
		} else {
			return n, done, err
		}
	}

	buff = buff[:read]

	if buff[0] == []byte(" ")[0] {
		return n, done, fmt.Errorf("Header cannot start with a whitespace: %s", buff)
	}

	indexOfSeparator := strings.Index(string(buff), ":")

	if indexOfSeparator == -1 {
		return n, done, fmt.Errorf("Header KVP's have to be separated by a ':'")
	}
	kvp := []string{string(buff[:indexOfSeparator]), string(buff[indexOfSeparator+1:])}

	keyHasTrailingWhitespace := strings.HasSuffix(kvp[0], " ")
	if keyHasTrailingWhitespace {
		return n, done, fmt.Errorf("Header key cannot have trailing white spaces")
	}
	key := strings.Fields(kvp[0])
	if len(key) != 1 {
		return n, done, fmt.Errorf("Header key cannot have white spaces in the middle")
	}

	valueTrimmed := strings.TrimSpace(kvp[1])
	h[key[0]] = valueTrimmed

	n = read + len([]byte(Clrf))
	return n, done, nil
}
