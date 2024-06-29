package main

import (
	"fmt"
	"strings"
)

type Http struct {
	Method     string
	Path       string
	Version    string
	StatusCode int
	Reason     string

	Header *Headers
}

type Headers struct {
	Host      string
	UserAgent string
	MediaType string
}

func New() *Http {
	return &Http{
		Version:    "HTTP/1.1",
		StatusCode: 200,
		Reason:     "OK",
	}
}

func ParseHttpRequest(data []byte) *Http {
	requestLine := []string{}
	// Search for spaces in first block
	index := 0
	for i := 0; i < len(data); i++ {
		if i == 20 { // Check for spaces
			fmt.Println(string(data[index:i]))
			requestLine = append(requestLine, string(data[index:i]))
			index = i
		}
	}

	requestLine = strings.Split(requestLine[0], " ")
	statusCode := 200
	reason := "OK"

	if requestLine[1] != "/" {
		statusCode = 404
		reason = "Not Found"
	}

	return &Http{
		Method:     requestLine[0],
		Path:       requestLine[1],
		Version:    "HTTP/1.1",
		StatusCode: statusCode,
		Reason:     reason,
	}
}

func (h *Http) Response() []byte {
	// CRLF - Carriage Return Line Feed
	// Also known as 'control characters'
	// CR - Moves the cursor to the beginning of the line without advancing to the next
	// LF - Moves the cursor down to the next line without returning to the beginning of the line.
	str := fmt.Sprintf("%s %d %s\r\n\r\n", h.Version, h.StatusCode, h.Reason)
	return []byte(str)
}
