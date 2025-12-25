package config

import "time"

type Config struct {
	ReadBufferSize int
	BodyLimit      int
	HeaderLimit    int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func Load(
	ReadBufferSize, BodyLimit, HeaderLimit int,
	ReadTimeout, WriteTimeout time.Duration) *Config {
	return &Config{
		ReadBufferSize: ReadBufferSize,
		BodyLimit:      BodyLimit,
		HeaderLimit:    HeaderLimit,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
	}
}
