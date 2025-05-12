package redis

import (
	"my-clean-architecture-template/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Connection struct {
	Config
	Connection *redis.Client
	Stop       chan struct{}

	Logger logger.Interface
}

// New -.
// func New(url string, opts ...Option) (*Connection, error) {
// 	pg := &Connection{
// 		maxPoolSize:  _defaultMaxPoolSize,
// 		connAttempts: _defaultConnAttempts,
// 		connTimeout:  _defaultConnTimeout,
// 	}

// 	// Custom options
// 	for _, opt := range opts {
// 		opt(pg)
// 	}

// 	return pg, nil
// }
