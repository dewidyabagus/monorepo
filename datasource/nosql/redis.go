package nosql

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"-"`
	DB           int    `json:"db,omitempty"`
	UseTLS       bool   `json:"use_tls,omitempty"`
	MaxRetries   int    `json:"max_retries"`
	MinIdleConns int    `json:"min_idle_conns"`
	PoolSize     int    `json:"pool_size"`
	PoolTimeout  int    `json:"pool_timeout"`
	MaxConnAge   int    `json:"max_conn_age"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// RedisNewClient open redis session with connection pooling, adjustment timeout and custome options
func RedisNewClient(config *Redis) (*redis.Client, error) {
	// Redis connection options
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		MinIdleConns: config.MinIdleConns,
		PoolSize:     config.PoolSize,
		PoolTimeout:  time.Second * time.Duration(config.PoolTimeout),  // Seconds
		MaxConnAge:   time.Second * time.Duration(config.MaxConnAge),   // Seconds
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),  // Seconds
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout), // Seconds
	}
	if config.UseTLS {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Open New Session
	rdb := redis.NewClient(options)

	// Test Connection And Auth with PING
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping command: %w", err)
	}

	return rdb, nil
}
