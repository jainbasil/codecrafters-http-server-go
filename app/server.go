package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Printf("Listening on host: localhost, port: 4221\n")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatal("Error closing the listener")
		}
	}(l)

	for {
		conn, err := l.Accept()
		log.Println("Connection received")
		if err != nil {
			log.Fatal("Error accepting connection: ", err.Error())
		}

		// handle the connection
		go func(conn net.Conn) {
			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
			}
			req := string(bytes.NewBuffer(buf).String())
			fmt.Println(req)

			// split by \r\n and get the first line --> request line.
			requestLine := strings.Split(req, "\r\n")[0]
			requestTarget := strings.Fields(requestLine)[1]

			response := "HTTP/1.1 200 OK\r\n\r\n"

			// only supports /
			if requestTarget != "/" {
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			}

			conn.Write([]byte(response))
			conn.Close()
		}(conn)
	}
}
