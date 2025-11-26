package main

import (
	"bytes"
	"log"
	"net"
	"time"
)

const (
	MAX_REQUEST_SIZE = 2 * 1024 * 1024
	READ_BUFFER      = 4 * 1024
	READ_TIMEOUT     = time.Second * 10
	WRITE_TIMEOUT    = time.Second * 10
)

func main() {
	log.Print("Server started")
	listener, err := net.Listen("tcp", ":1783")
	if err != nil {
		log.Fatalf("Listening Socket Error : %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection Error : %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(READ_TIMEOUT))

	buffer := make([]byte, READ_BUFFER)      // No global default, set accordingly
	req := make([]byte, 0, MAX_REQUEST_SIZE) // Again Http doesnt provide this, depends on server logic

	for {
		streamLength, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Read Error :%v", err)
			return
		}

		log.Printf("Data : \n%v,StreamLength : %d\n", string(buffer[:streamLength]), streamLength)
		if len(req)+streamLength > MAX_REQUEST_SIZE {
			return
		}

		req = append(req, buffer[:streamLength]...)
		if bytes.Contains(req, []byte("\r\n\r\n")) {
			break
		}

	}
	log.Printf("Full Request : %s\n", string(req))

	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 5\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello"
	respByte := []byte(resp)

	conn.SetWriteDeadline(time.Now().Add(WRITE_TIMEOUT))

	_, err := conn.Write(respByte)
	if err != nil {
		log.Printf("Write Error :%v", err)
	}
}
