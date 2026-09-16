package headers

import (
	"bytes"
	"fmt"
	"slices"
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

	kvpLine := data[:endIndex]
	read := len(kvpLine)

	if kvpLine[0] == []byte(" ")[0] {
		return n, done, fmt.Errorf("Header cannot start with a whitespace: %s", kvpLine)
	}

	indexOfSeparator := strings.Index(string(kvpLine), ":")

	if indexOfSeparator == -1 {
		return n, done, fmt.Errorf("Header KVP's have to be separated by a ':'")
	}
	kvp := []string{string(kvpLine[:indexOfSeparator]), string(kvpLine[indexOfSeparator+1:])}

	keyHasTrailingWhitespace := strings.HasSuffix(kvp[0], " ")
	if keyHasTrailingWhitespace {
		return n, done, fmt.Errorf("Header key cannot have trailing white spaces")
	}
	keyParts := strings.Fields(kvp[0])
	if len(keyParts) != 1 {
		return n, done, fmt.Errorf("Header key cannot have white spaces in the middle")
	}
	key := strings.ToLower(keyParts[0])
	validKey := true
	for _, v := range key {
		if (v >= 'A' && v <= 'Z') ||
			(v >= 'a' && v <= 'z') ||
			(v >= '0' && v <= '9') ||
			slices.Contains(tokenChars, v) {
			continue
		} else {
			validKey = false
			break
		}
	}

	if !validKey {
		return n, done, fmt.Errorf("Header key cannot has invalid characters")
	}
	valueTrimmed := strings.TrimSpace(kvp[1])

	existingVal, ok := h[key]
	if ok {
		h[key] = fmt.Sprintf("%s, %s", existingVal, valueTrimmed)
	} else {
		h[key] = valueTrimmed
	}

	n = read + len([]byte(Clrf))
	return n, done, nil
}

var tokenChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
