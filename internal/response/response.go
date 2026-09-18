package response

import (
	"fmt"
	"httpz/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusCodeOk                  = 200
	StatusCodeBadRequest          = 400
	StatusCodeInternalServerError = 500
)

const clrf = "\r\n"

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusCodeOk:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	responseBytes := []byte("HTTP/1.1 ")

	switch statusCode {
	case StatusCodeOk:
		responseBytes = append(responseBytes, []byte("200 OK")...)
	case StatusCodeBadRequest:
		responseBytes = append(responseBytes, []byte("400 Bad Request")...)
	case StatusCodeInternalServerError:
		responseBytes = append(responseBytes, []byte("500 Internal Server Error")...)
	default:
		responseBytes = append(responseBytes, []byte(strconv.Itoa(int(statusCode)))...)
	}

	responseBytes = append(responseBytes, []byte(clrf)...)

	_, err := w.Write(responseBytes)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders["Content-Length"] = strconv.Itoa(contentLen)
	defaultHeaders["Connection"] = "close"
	defaultHeaders["Content-Type"] = "text/plain"
	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headersBytes := make([]byte, 0)
	for key, value := range headers {
		headerString := fmt.Sprintf("%s: %s%s", key, value, clrf)
		headersBytes = append(headersBytes, []byte(headerString)...)
	}

	headersBytes = append(headersBytes, []byte(clrf)...) // end of headers

	_, err := w.Write(headersBytes)
	return err
}
