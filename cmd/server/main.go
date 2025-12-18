package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
	"github.com/brutally-Honest/http-server/internal/request"
	"github.com/brutally-Honest/http-server/internal/response"
	"github.com/brutally-Honest/http-server/internal/router"
	"github.com/brutally-Honest/http-server/internal/server"
)

const (
	DefaultBufferSize = 4 * 1024
	MaxBodySize       = 2 * 1024 * 1024
	MaxHeaderSize     = 8 * 1024
	ReadTimeout       = time.Second * 10
	WriteTimeout      = time.Second * 10
)

func main() {
	cfg := config.Load(
		DefaultBufferSize,
		MaxBodySize,
		MaxHeaderSize,
		ReadTimeout,
		WriteTimeout,
	)

	r := router.NewRouter()
	r.GET("/api/static", func(req *request.Request, res *response.Response) {
		res.Write([]byte("WOHOO !!! It is working"))
		res.Flush(req, false)
	})
	r.GET("/api/param/:id", func(req *request.Request, res *response.Response) {
		id := req.Params["id"]
		output := fmt.Sprintf("Id %s", id)
		res.Write([]byte(output))
		res.Flush(req, false)
	})
	r.GET("/api/param/:id/profile/:name", func(req *request.Request, res *response.Response) {
		id := req.Params["id"]
		name := req.Params["name"]

		output := fmt.Sprintf("Id %s Name %s", id, name)
		res.Write([]byte(output))
		res.Flush(req, false)
	})
	r.GET("/api/wildcard/*anything", func(req *request.Request, res *response.Response) {
		wildcard := req.Params["anything"]

		output := fmt.Sprintf("wild path %s", wildcard)
		res.Write([]byte(output))
		res.Flush(req, false)
	})
	r.POST("/api/wake-up", func(req *request.Request, res *response.Response) {
		log.Print(string(req.Body))
		//simulating something created
		res.WriteHeader(201)
		res.Flush(req, false)
	})

	r.GET("/stream", func(req *request.Request, res *response.Response) {
		res.SetHeader("Content-Type", "text/plain")
		res.SetHeader("Transfer-Encoding", "chunked")

		chunks := []string{"Testing\n", "Transfer\n", "Encoding\n", "With\n", "HTTP\n", "1.1\n"}

		for _, chunk := range chunks {
			if err := res.WriteChunk([]byte(chunk)); err != nil {
				log.Printf("WriteChunk failed: %v", err)
				return
			}
		}

		if err := res.EndChunked(); err != nil {
			log.Printf("EndChunked failed: %v", err)
			return
		}
	})

	s := server.NewServer(":1783", cfg, r)
	log.Fatal(s.ListenAndServe())
}
