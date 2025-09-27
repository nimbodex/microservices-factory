package app

import (
	"github.com/nimbodex/microservices-factory/payment/internal/config"
)

// Container представляет контейнер зависимостей
type Container struct {
	config *config.Config
}

// NewContainer создает новый контейнер зависимостей
func NewContainer() *Container {
	return &Container{
		config: config.New(),
	}
}

func (c *Container) Config() *config.Config {
	return c.config
}

func (c *Container) BuildApp() *App {
	return New(c.config)
}
