package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var (
	wg       sync.WaitGroup
	filesDir string
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Check for --directory arg
	if len(os.Args) > 1 && os.Args[1] == "--directory" {
		filesDir = os.Args[2]
	} else {
		filesDir = ""
	}

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Could not read from connection: %v", err)
	}

	h := NewParseHttpRequest(buf, filesDir)
	resBytes := h.Response()

	_, err = conn.Write(resBytes)
	if err != nil {
		log.Fatal(err)
	}

	return
}
