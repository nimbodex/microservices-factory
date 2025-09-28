package mongo

import (
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

type Config struct {
	Image       string
	Port        string
	Database    string
	Username    string
	Password    string
	StartupWait wait.Strategy
}

func DefaultConfig() *Config {
	return &Config{
		Image:    "mongo:7.0",
		Port:     "27017/tcp",
		Database: "test",
		Username: "test",
		Password: "test",
		StartupWait: wait.ForListeningPort("27017/tcp").
			WithStartupTimeout(60 * time.Second).
			WithPollInterval(1 * time.Second),
	}
}
