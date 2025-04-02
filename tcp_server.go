package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Server interface {
	Run() error
}

type Tcp_Server struct {
	port string
}

type Client struct {
	rd *bufio.Reader
	// wt  *bufio.Writer
	con net.Conn
}

// Possible typ are
// Simple strings 	RESP2 	Simple 	+
// Integers 	RESP2 	Simple 	:
// Bulk strings 	RESP2 	Aggregate 	$
type RESPValue struct {
	typ     string
	integer int
	str     string
	arr     []RESPValue
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
	// defer conn.Close()
	client := Client{
		rd: bufio.NewReader(conn),
		// wt:  bufio.NewWriter(conn),
		con: conn,
	}
	for {

		RESPInput := readVal(client.rd)

		fmt.Println("handleConn.readVal RESPInput=", RESPInput)

		bulkArray := RESPInput.arr
		fmt.Println("handleConn.bulkArray=", bulkArray)
		// 		At this moment assuming a null command wont be sent.
		// 		Redis commands are all uppercase, so I will enforce the same.
		command := strings.ToUpper(bulkArray[0].str)

		fmt.Println("handleConn =", command)

		switch command {
		case "PING":

			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				fmt.Println("Write Response Error:", err)
				errCh <- err

				// return err
			}
		case "ECHO":
			msg := bulkArray[1].str
			_, err := conn.Write([]byte("+'" + msg + "'\r\n"))
			if err != nil {
				fmt.Println("Write Response Error:", err)
				errCh <- err

				// return err
			}
		default:
			_, err := conn.Write([]byte("\r\n"))
			if err != nil {
				fmt.Println("Write Response Error:", err)
				errCh <- err
			}
		}
	}
}

func readVal(rd *bufio.Reader) (vals RESPValue) {

	var bt byte

	bt, err := rd.ReadByte()
	if err != nil {
		fmt.Println("readVal.ReadByte Error: ", err)

	}

	if bt == '*' {
		// Now we know we have a multibulk array
		// Next step is to get the length and proceed with reading it
		// vals = append(vals, readArrayVal(rd))
		vals = readArrayVal(rd)

	} else {

		// Simple strings 	RESP2 	Simple 	+
		if bt == '+' {

		}
		// Integers 	RESP2 	Simple 	:
		if bt == ':' {

		}
		// Bulk strings 	RESP2 	Aggregate 	$
		if bt == '$' {
			vals = readBulkVal(rd)

		}
	}
	return vals
}

func readArrayVal(rd *bufio.Reader) (val RESPValue) {
	btint, err := rd.ReadBytes('\n')
	if err != nil {
		fmt.Println("readArrayVal.ReadBytes Error: ", err)

	}
	// ReadBytes will also read \r\n so [:len(btint)-2] will remove the
	// last 2 bytes which wil be \r\n
	bulkArrLen, err := strconv.ParseInt(string(btint[:len(btint)-2]), 10, 64)
	if err != nil {
		fmt.Println("readArrayVal.ParseInt Error: ", err)

	}
	fmt.Println("readArrayVal bulkArrLen=", bulkArrLen)

	// In the case of nested arrays, bulkArrLen is doesnt count number of items
	// in the nested array. So i thinkk the array will be atleast bulkArrLen and will
	// have to dynamically resize itself
	// Actually atm I dont think nested arrays will even parse
	bulkArrayVals := make([]RESPValue, bulkArrLen)

	for i := 0; i < int(bulkArrLen); i++ {
		// Now we need to read each element in the bulk array
		// There are bulkArrLen amount of elements

		val := readVal(rd)
		bulkArrayVals[i] = val
	}
	return RESPValue{typ: "*", arr: bulkArrayVals}
}

func readBulkVal(rd *bufio.Reader) (val RESPValue) {

	btint, err := rd.ReadBytes('\n')
	if err != nil {
		fmt.Println("readVal.$.ReadBytes Error: ", err)

	}
	bulkValLen, err := strconv.ParseInt(string(btint[:len(btint)-2]), 10, 64)
	if err != nil {
		fmt.Println("readVal.$.ParseInt Error: ", err)

	}
	fmt.Println("readVal bulkValLen=", bulkValLen)

	//	+2 because \n\r will also be read
	bulkStringBuff := make([]byte, bulkValLen+2)
	nread, err := io.ReadFull(rd, bulkStringBuff)
	if err != nil {
		fmt.Println("readVal.$.ReadFull Error: ", err)
	}
	fmt.Println("readVal.$.ReadFull nread=", nread)
	s := fmt.Sprintf("command:\n%s", bulkStringBuff)
	fmt.Println(s)
	return RESPValue{typ: "$", str: strings.TrimRight(string(bulkStringBuff), "\r\n")}

}
