package main

import (
	"bytes"
	"log"
	"net"
)

const MAX_REQ_SIZE = 1024

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1783")
	if err != nil {
		log.Fatalf("Listening Socket Error : %v", err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Connection Error : %v", err)
	}
	log.Println("Remote:", conn.RemoteAddr())
	log.Println("Local:", conn.LocalAddr())

	buffer := make([]byte, 80)
	req := make([]byte, 0, MAX_REQ_SIZE)
	for {
		streamLength, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("Internal Error :%v", err)
		}
		log.Printf("Data :%v , StreamLength : %d", string(buffer[:streamLength]), streamLength)
		if len(req)+streamLength > MAX_REQ_SIZE {
			conn.Close()
			return
		}
		req = append(req, buffer[:streamLength]...)
		if bytes.Contains(req, []byte("\r\n\r\n")) {
			break
		}

	}
	log.Print("Request", string(req))
	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 5\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello"
	respByte := []byte(resp)
	log.Print("Response BYtes Length:", len(respByte))
	n, err := conn.Write(respByte)
	log.Print("Written Data Length :", n)
	conn.Close()
}
