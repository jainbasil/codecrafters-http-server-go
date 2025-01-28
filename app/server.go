package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var FileDirectory string

func init() {
	flag.StringVar(&FileDirectory, "directory", ".", "the directory to serve files from")
	flag.Parse()
}

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
			req := ParseRequest(buf)
			res := HandleRequest(req)

			conn.Write(SerializeReponse(res))
			conn.Close()
		}(conn)
	}
}
