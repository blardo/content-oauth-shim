package config

import (
	"fmt"

	redistore "gopkg.in/boj/redistore.v1"
)

// ServerOption is an option type for ServerConfig
type ServerOption func(*ServerConfig)

// RedisConnection sets the redis connection string for the server configuration
func RedisConnection(hostport string) ServerOption {
	return func(sc *ServerConfig) {
		sc.redisConnection = hostport
	}
}

// Port sets the port on the ServerConfig
func Port(port int) ServerOption {
	return func(sc *ServerConfig) {
		sc.Port = port
	}
}

// NewServerConfig returns a server configuration with the defaults
//
// 		Host: "0.0.0.0"
//		Port: 3000
// 		redisConnection: "192.168.99.100:6379"
//
func NewServerConfig(sessionSecret string, options ...ServerOption) (*ServerConfig, error) {
	sc := &ServerConfig{
		Host:            "0.0.0.0",
		Port:            3000,
		sessionSecret:   []byte(sessionSecret),
		redisConnection: "192.168.99.100:6379",
	}
	for i := range options {
		if options[i] != nil {
			options[i](sc)
		}
	}

	rStore, err := redistore.NewRediStore(10, "tcp", sc.redisConnection, "", sc.sessionSecret)
	if err != nil {
		return nil, err
	}
	sc.RedisStore = rStore
	return sc, nil
}

// ServerConfig is a configuration struct for the server
type ServerConfig struct {
	Host       string               // The address to listen on
	Port       int                  // The port to listen on
	RedisStore *redistore.RediStore // The redis connection

	sessionSecret   []byte // A secret for sessions
	redisConnection string // The host:port to connect to redis
}

// Hostport returns a listen string for http.Listen()
func (s ServerConfig) Hostport() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
