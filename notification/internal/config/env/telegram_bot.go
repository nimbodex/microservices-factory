package env

import (
	"os"
)

type TelegramBotConfig struct{}

func NewTelegramBotConfig() *TelegramBotConfig {
	return &TelegramBotConfig{}
}

func (c *TelegramBotConfig) GetBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func (c *TelegramBotConfig) GetChatID() string {
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		return "123456789" // Захардкоженный ID для демонстрации
	}
	return chatID
}
