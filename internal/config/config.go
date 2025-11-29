package config

import "time"

type Config struct {
	BufferLimit  int
	BodyLimit    int
	HeaderLimit  int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load(
	BufferLimit, BodyLimit, HeaderLimit int,
	ReadTimeout, WriteTimeout time.Duration) *Config {
	return &Config{
		BufferLimit:  BufferLimit,
		BodyLimit:    BodyLimit,
		HeaderLimit:  HeaderLimit,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}
}
