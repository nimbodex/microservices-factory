package mongo

import "github.com/testcontainers/testcontainers-go/wait"

type Option func(*Config)

func WithImage(image string) Option {
	return func(c *Config) {
		c.Image = image
	}
}

func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

func WithCredentials(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

func WithStartupWait(wait wait.Strategy) Option {
	return func(c *Config) {
		c.StartupWait = wait
	}
}
