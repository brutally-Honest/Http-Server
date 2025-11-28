package config

import "time"

type Config struct {
	BufferLimit  int
	RequestLimit int
	HeaderLimit  int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load(
	BufferLimit, RequestLimit, HeaderLimit int,
	ReadTimeout, WriteTimeout time.Duration) *Config {
	return &Config{
		BufferLimit:  BufferLimit,
		RequestLimit: RequestLimit,
		HeaderLimit:  HeaderLimit,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}
}
