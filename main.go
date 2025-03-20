package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer ln.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Accept incoming connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	fmt.Println("New client")

	conn.Write([]byte("Hello from server"))

	defer conn.Close()
}
