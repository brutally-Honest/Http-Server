package main

import (
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
	"github.com/brutally-Honest/http-server/internal/server"
)

func main() {
	cfg := config.Load(
		4*1024,
		2*1024*1024,
		8*1024,
		time.Second*10,
		time.Second*10,
	)
	s := server.NewServer(":1783", cfg)
	s.ListenAndServe()
}
