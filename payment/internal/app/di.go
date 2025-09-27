package app

import (
	"github.com/maxim/microservices-factory/payment/internal/config"
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

// Config возвращает конфигурацию
func (c *Container) Config() *config.Config {
	return c.config
}

// BuildApp создает и настраивает приложение
func (c *Container) BuildApp() *App {
	return New(c.config)
}
