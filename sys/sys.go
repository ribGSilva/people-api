package sys

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Configs contains all the configs gathered from env vars
var Configs struct {
	IdempotencyEnabled bool
	Http               struct {
		Port            string
		ShutdownTimeout time.Duration
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
	}
	Swagger struct {
		Protocol string
		Host     string
	}
	Mongo struct {
		Database          string
		ConnectionURL     string
		ConnectionTimeout time.Duration
		DisconnectTimeout time.Duration
		PingTimeout       time.Duration
		OperationTimeout  time.Duration
	}
	Redis struct {
		ConnectionURL    string
		User             string
		Pass             string
		PingTimeout      time.Duration
		OperationTimeout time.Duration
		CacheTTL         time.Duration
	}
	NewRelic struct {
		AppName           string
		Licence           string
		Enabled           bool
		ConnectionTimeout time.Duration
		ShutdownTimeout   time.Duration
	}
}

// S holds static resources across the project
var S struct {
	Log   *zap.SugaredLogger
	Mongo *mongo.Client
	Redis *redis.Client
}
