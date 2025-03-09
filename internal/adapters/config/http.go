package config

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net"
)

// HTTPConfig defines an interface for HTTP server configuration.
type HTTPConfig interface {
	Address() string
	EnabledTLS() bool
	CertFile() string
	KeyFile() string
	GetCORSConfig() middleware.CORSConfig
}

type httpConfig struct {
	host         string
	port         string
	enabledTLS   bool
	tlsPort      string
	certFile     string
	keyFile      string
	allowOrigins []string
}

// NewHTTPConfig initializes a new HTTP configuration from environment variables.
func NewHTTPConfig() (HTTPConfig, error) {
	var allowOrigins []string
	if viper.GetBool("backend.dev-mode") {
		allowOrigins = viper.GetStringSlice("backend.cors.dev.allow-origins")
	} else {
		allowOrigins = viper.GetStringSlice("backend.cors.prod.allow-origins")
	}

	return &httpConfig{
		host:         viper.GetString("backend.host"),
		port:         viper.GetString("backend.port"),
		enabledTLS:   viper.GetBool("backend.tls.enabled"),
		tlsPort:      viper.GetString("backend.tls.port"),
		certFile:     viper.GetString("backend.tls.cert-file"),
		keyFile:      viper.GetString("backend.tls.key-file"),
		allowOrigins: allowOrigins,
	}, nil
}

// Address constructs and returns the full server address (host:port).
func (cfg *httpConfig) Address() string {
	if cfg.enabledTLS {
		return net.JoinHostPort(cfg.host, cfg.tlsPort)
	}
	return net.JoinHostPort(cfg.host, cfg.port)
}

// EnabledTLS returns true if TLS is enabled.
func (cfg *httpConfig) EnabledTLS() bool {
	return cfg.enabledTLS
}

// CertFile returns the path to the TLS certificate file.
func (cfg *httpConfig) CertFile() string {
	return cfg.certFile
}

// KeyFile returns the path to the TLS key file.
func (cfg *httpConfig) KeyFile() string {
	return cfg.keyFile
}

func (cfg *httpConfig) GetCORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     cfg.allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		MaxAge:           300, // Maximum age for browser to cache preflight request results
	}
}
