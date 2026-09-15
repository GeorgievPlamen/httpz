package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error reading all bytes: %v", err)
	}

	requestLine, err := parseRequestLine(bytes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing requestline: %v", err)
	}

	return &Request{RequestLine: *requestLine}, nil
}

const httpNewLines = "\r\n"

func parseRequestLine(bytes []byte) (*RequestLine, error) {
	parts := strings.Split(string(bytes), httpNewLines)
	requestLineRaw := parts[0]

	requestLineParts := strings.Fields(requestLineRaw)
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("Invalid number of request line parts, must be 3")
	}

	method := requestLineParts[0]
	path := requestLineParts[1]
	version := strings.Split(requestLineParts[2], "/")

	for _, v := range method {
		isCapitalLetter := unicode.IsLetter(v) && unicode.IsUpper(v)
		if !isCapitalLetter {
			return nil, fmt.Errorf("Method invalid, must be only capital letters: %v", method)
		}
	}

	if len(version) != 2 || version[1] != "1.1" {
		return nil, fmt.Errorf("Invalid HTTP version: %v", version)
	}

	return &RequestLine{
		HttpVersion:   version[1],
		RequestTarget: path,
		Method:        method,
	}, nil
}
