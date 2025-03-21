package main

import (
	"fmt"
	"io"
	"net"
)

type Server interface {
	Run() error
}

type Tcp_Server struct {
	port string
}

func (t Tcp_Server) Run() error {

	ln, err := net.Listen("tcp", t.port)
	if err != nil {
		fmt.Println("Listen Error:", err)
		return err
	}
	defer ln.Close()

	fmt.Println("Server is listening on port ", t.port)
	for {
		// Accept incoming connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			return err
		}

		errCh := make(chan error)
		go handleConn(conn, errCh)

		// errs := <-errCh
		// if errs != nil {
		// 	return errs
		// }
	}
	return nil
}

func handleConn(conn net.Conn, errCh chan<- error) {
	defer conn.Close()

	buff := make([]byte, 128)
	for {
		nread, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				errCh <- nil
				// return nil
			}

			fmt.Println("Read Error: ", err)
			errCh <- err
			// return err
		}
		fmt.Printf("command:\n%s", buff[:nread])

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Write Response Error:", err)
			errCh <- err

			// return err
		}
	}
}
