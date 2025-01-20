package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Printf("Listening on host: localhost, port: 4221\n")

	// Uncomment this block to pass the first stage
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
			fmt.Println(bytes.NewBuffer(buf).String())

			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			conn.Close()
		}(conn)
	}
}
