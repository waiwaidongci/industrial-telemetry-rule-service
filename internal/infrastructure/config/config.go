package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	ShutdownSeconds int
	ShutdownTimeout time.Duration
	Environment     string
}

func Load() Config {
	c := Config{HTTPAddr: env("HTTP_ADDR", ":8086"), ShutdownSeconds: 10, ShutdownTimeout: LoadDuration("SHUTDOWN_TIMEOUT", 10*time.Second), Environment: env("APP_ENV", "development")}
	if n, err := strconv.Atoi(os.Getenv("SHUTDOWN_SECONDS")); err == nil {
		c.ShutdownSeconds = n
	}
	return c
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
