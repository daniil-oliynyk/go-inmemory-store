package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Server interface {
	Run() error
}

type Tcp_Server struct {
	port string
}

type Client struct {
	rd  *bufio.Reader
	wt  *bufio.Writer
	con net.Conn
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

	client := Client{
		rd:  bufio.NewReader(conn),
		wt:  bufio.NewWriter(conn),
		con: conn,
	}

	// buff := make([]byte, 128)
	for {

		var bt byte

		bt, err := client.rd.ReadByte()
		if err != nil {
			fmt.Println("ReadByte Error: ", err)
			errCh <- err
		}

		if bt == '*' {
			// Now we know we have a multibulk array
			// Next step is to get the length and proceed with reading it
			btint, err := client.rd.ReadBytes('\n')
			if err != nil {
				fmt.Println("ReadBytes Error: ", err)
				errCh <- err
			}
			// ReadBytes will also read \r\n so [:len(btint)-2] will remove the
			// last 2 bytes which wil be \r\n
			bulkArrLen, err := strconv.ParseInt(string(btint[:len(btint)-2]), 10, 64)
			if err != nil {
				fmt.Println("ParseInt Error: ", err)
				errCh <- err
			}
			fmt.Println(bulkArrLen)

			for i := 0; i < int(bulkArrLen); i++ {
				// Now we need to read each element in the bulk array
				// There are bulkArrLen amount of elements

			}

		}

		// nread, err := conn.Read(buff)
		// if err != nil {
		// 	if err == io.EOF {
		// 		errCh <- nil
		// 		// return nil
		// 	}

		// 	fmt.Println("Read Error: ", err)
		// 	errCh <- err
		// 	// return err
		// }
		// s := fmt.Sprintf("command:\n%s", buff[:nread])
		// fmt.Println(s)

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Write Response Error:", err)
			errCh <- err

			// return err
		}
	}
}
