package config

import (
	"fmt"
	"net"
)

func (v Config) Validate() error {
	if v.HTTPAddr == "" {
		return fmt.Errorf("http address required")
	}
	if _, _, err := net.SplitHostPort(v.HTTPAddr); err != nil {
		return fmt.Errorf("invalid http address: %w", err)
	}
	if v.ShutdownSeconds <= 0 {
		return fmt.Errorf("shutdown seconds must be positive")
	}
	if v.ShutdownTimeout == 0 {
		return fmt.Errorf("shutdown timeout must be configured")
	}
	return nil
}
func (v Config) Redacted() map[string]any {
	return map[string]any{"http_addr": v.HTTPAddr, "environment": v.Environment, "shutdown_seconds": v.ShutdownSeconds}
}
