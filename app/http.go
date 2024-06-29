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

	Header Headers

	ResponseBody string
}

type Headers struct {
	Host          string
	UserAgent     string
	ContentType   string
	ContentLength int
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
	headers := []string{}
	// Search for spaces in first block
	index := 0
	for i := 0; i < len(data); i++ {
		if i == 20 { // Check for spaces
			requestLine = append(requestLine, string(data[index:i]))
			index = i
		}

		if i == 10 {
			headers = append(headers, string(data[i:]))
		}
	}

	requestLine = strings.Split(requestLine[0], " ")
	statusCode := 200
	reason := "OK"
	h := &Headers{}
	resBody := ""

	if strings.Contains(requestLine[1], "/echo") {
		str := strings.Split(requestLine[1], "/")
		resBody = str[2]
		h.ContentType = "text/plain"
		h.ContentLength = len(str[2])
	} else if requestLine[1] != "/" {
		statusCode = 404
		reason = "Not Found"
	}

	return &Http{
		Method:       requestLine[0],
		Path:         requestLine[1],
		Version:      "HTTP/1.1",
		StatusCode:   statusCode,
		Reason:       reason,
		Header:       *h,
		ResponseBody: resBody,
	}
}

func (h *Http) Response() []byte {
	// CRLF - Carriage Return Line Feed
	// Also known as 'control characters'
	// CR - Moves the cursor to the beginning of the line without advancing to the next
	// LF - Moves the cursor down to the next line without returning to the beginning of the line.
	str := fmt.Sprintf("%s %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s\r\n\r\n", h.Version, h.StatusCode, h.Reason, h.Header.ContentType, h.Header.ContentLength, h.ResponseBody)
	return []byte(str)
}
