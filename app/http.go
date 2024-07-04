package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	OK          = 200
	NOT_FOUND   = 404
	BAD_REQUEST = 400
)

type Http struct {
	RequestLine  HttpRequestLine
	Headers      HttpHeaders
	ResponseBody HttpResponse
}

type HttpRequestLine struct {
	Method     string
	Path       string
	Version    string
	Reason     string
	StatusCode int
}

type HttpHeaders struct {
	Host          string
	UserAgent     string
	Accept        string
	ContentType   string
	ContentLength int
}

type HttpResponse struct {
	Body string
}

// https://datatracker.ietf.org/doc/html/rfc9112#name-message-parsing
func NewParseHttpRequest(request []byte) *Http {
	// Http request is broke down as follows
	/*
					   Request Line: Seperated by spaces, and ends with a control character
					     - Method
					     - Path
					     - Version

				    Headers: Each header is seperated with a control character.
				      - Host
				      - User-Agent
				      - Content-Type
				      - Etc...

		        ResponseBody: Optional response body
		          - ...

		        Each request should end with a double control character /r/n/r/n
	*/
	var blocks [][]byte
	crlf := []byte{13, 10}

	for len(request) > 0 {
		// Grab the index of a CRLF block
		index := bytes.Index(request, crlf)

		// No more CRLF blocks found
		if index == -1 {
			blocks = append(blocks, request)
			break
		}

		// Add block up to and including the CRLF
		block := request[:index+len(crlf)]
		blocks = append(blocks, block)

		// Move past CRLF
		request = request[index+len(crlf):]
	}

	requestLine := parseLineRequest(blocks[0])
	headers, _ := parseHeaderRequest(blocks[1:])

	http := &Http{
		RequestLine: *requestLine,
		Headers:     *headers,
	}
	http.parsePathResponse(requestLine.Path)

	return http
}

func parseLineRequest(block []byte) *HttpRequestLine {
	data := strings.Split(string(block), " ")
	fmt.Println(data)
	if len(data) <= 1 {
		log.Fatal("Could not parse line request block.")
	}

	return &HttpRequestLine{
		Method:     data[0],
		Path:       strings.TrimSpace(data[1]),
		Version:    "HTTP/1.1",
		StatusCode: OK,
		Reason:     "OK",
	}
}

func (h *Http) parsePathResponse(path string) {
	switch path {
	case "/":
		h.RequestLine.Reason = "OK"
		h.RequestLine.StatusCode = OK
	default:
		if strings.Contains(path, "/echo") {
			str := strings.Split(path, "/")
			h.Headers.ContentType = "text/plain"
			h.Headers.ContentLength = len(str[len(str)-1])
			h.RequestLine.StatusCode = OK
			h.ResponseBody.Body = str[len(str)-1]
			return
		} else if strings.Contains(path, "/user-agent") {
			h.ResponseBody.Body = h.Headers.UserAgent
			h.Headers.ContentType = "text/plain"
			h.Headers.ContentLength = len(h.Headers.UserAgent)
			h.RequestLine.StatusCode = OK
			return
		} else {
			h.RequestLine.Reason = "Not Found"
			h.RequestLine.StatusCode = NOT_FOUND
		}
	}
}

func parseHeaderRequest(headerBlocks [][]byte) (*HttpHeaders, int) {
	headers := map[string]string{}
	crlf := []byte{13, 10}
	index := 1

	for _, v := range headerBlocks {
		if string(v) == string(crlf) {
			break
		}

		h := strings.Split(string(v), ":")
		headers[h[0]] = h[1]
		index++
	}

	h := &HttpHeaders{}

	for key, val := range headers {
		switch key {
		case "Host":
			h.Host = strings.TrimSpace(val)
		case "User-Agent":
			h.UserAgent = strings.TrimSpace(val)
		case "Accept":
			h.Accept = strings.TrimSpace(val)
		case "Content-Type":
			h.ContentType = strings.TrimSpace(val)
		case "Content-Length":
			l, _ := strconv.Atoi(strings.TrimSpace(val))
			h.ContentLength = l
		}
	}

	if len(headerBlocks) < index {
		index = 0
	}

	// Return the current headers, as well as the index to the response block, or 0 if it doesnt exist
	return h, index
}

func (h *Http) parseBodyResponse(block [][]byte) {
	// TODO: Remember to deal with POST requests here when we come to it
}

func (h *Http) Response() []byte {
	// CRLF - Carriage Return Line Feed
	// Also known as 'control characters'
	// CR - Moves the cursor to the beginning of the line without advancing to the next
	// LF - Moves the cursor down to the next line without returning to the beginning of the line.
	str := fmt.Sprintf(
		"%s %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\nUser-Agent: %s\r\n\r\n%s",
		h.RequestLine.Version,
		h.RequestLine.StatusCode,
		h.RequestLine.Reason,
		h.Headers.ContentType,
		h.Headers.ContentLength,
		h.Headers.UserAgent,
		h.ResponseBody.Body,
	)

	return []byte(str)
}
