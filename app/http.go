package main

import "fmt"

type Http struct {
	Version    string
	StatusCode int
	Reason     string
}

func New() *Http {
	return &Http{
		Version:    "HTTP/1.1",
		StatusCode: 200,
		Reason:     "OK",
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
