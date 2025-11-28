package parser

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
)

func ParseRequest(conn net.Conn, cfg *config.Config) {

	buffer := make([]byte, cfg.BufferLimit)  // No global default, set accordingly
	req := make([]byte, 0, cfg.RequestLimit) // Again Http doesnt provide this, depends on server logic

	var headerIdx = -1
	conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	for {
		streamLength, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Read Error :%v", err)
			return
		}

		if len(req)+streamLength > cfg.HeaderLimit {
			log.Printf("Maximum Request size limit breached")
			return
		}

		req = append(req, buffer[:streamLength]...)
		if bytes.Contains(req, []byte("\r\n\r\n")) {
			headerIdx = bytes.Index(req, []byte("\r\n\r\n"))
			break
		}
	}

	headers := req[:headerIdx]
	contentLength, p_err := parseContentLength(headers)
	if p_err != nil {
		log.Print(p_err)
	}
	log.Println("Content Length:", contentLength)

	body := req[headerIdx+4:]
	remainingBody := contentLength - len(body)
	log.Printf("Remaining Body :%d", remainingBody)

	for remainingBody > 0 {
		bodyStreamLength, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Body Read Error :%v", err)
			return
		}
		log.Printf("Body Chunk: %d bytes", bodyStreamLength)

		if len(req)+bodyStreamLength > cfg.RequestLimit {
			log.Printf("Maximum Request size limit breached")
			return
		}

		req = append(req, buffer[:bodyStreamLength]...)
		body = append(body, buffer[:bodyStreamLength]...)

		remainingBody -= bodyStreamLength
	}

	log.Printf("Headers :%v", len(headers)+4)
	log.Printf("Body :%v", len(body))
	log.Printf("Request Length: %d\n", len(req))

}

func parseContentLength(headers []byte) (int, error) {
	lines := bytes.Split(headers, []byte("\r\n"))
	found := false
	var length int

	for _, line := range lines {
		if len(line) < 15 {
			continue
		}

		// Use EqualFold - case-insensitive without allocation
		if !bytes.EqualFold(line[:15], []byte("content-length:")) {
			continue
		}

		// Multiple Content-Length headers = error
		if found {
			return 0, fmt.Errorf("multiple Content-Length headers")
		}

		value := bytes.TrimSpace(line[15:])

		n, err := strconv.Atoi(string(value))
		if err != nil {
			return 0, fmt.Errorf("invalid Content-Length value: %w", err)
		}

		if n < 0 {
			return 0, fmt.Errorf("negative Content-Length: %d", n)
		}

		length = n
		found = true
	}

	return length, nil
}
