package app

import (
	"github.com/nimbodex/microservices-factory/inventory/internal/config"
)

type Container struct {
	config *config.Config
}

func NewContainer() *Container {
	return &Container{
		config: config.New(),
	}
}

func (c *Container) BuildApp() *App {
	return New(c.config)
}
