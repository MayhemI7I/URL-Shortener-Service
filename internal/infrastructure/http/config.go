package http

import (
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
)

type HTTPConfig struct {
	config.BaseConfig
	Address string
	Port string
	Protocol string // http or https
	ReadTimeout time.Duration
	WriteTimeout time.Duration

}

func NewHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		BaseConfig: config.NewBaseConfig("HTTP"),
		Address: "localhost",
		Port: "8080",
		Protocol: "http",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (c *HTTPConfig) AddFlags() {
	c.AddStringFlag(c.GetFlags(), "address", c.Address, "Address to listen on")
	c.AddStringFlag(c.GetFlags(), "port", c.Port, "Port to listen on")
	c.AddStringFlag(c.GetFlags(), "protocol", c.Protocol, "Protocol to listen on")
	c.AddDurationFlag(c.GetFlags(), "read-timeout", c.ReadTimeout, "Read timeout")
	c.AddDurationFlag(c.GetFlags(), "write-timeout", c.WriteTimeout, "Write timeout")
}


