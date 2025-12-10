package main

import (
	"log"
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
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
	s := server.NewServer(":1783", cfg, r)
	log.Fatal(s.ListenAndServe())
}
