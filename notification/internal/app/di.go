package app

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/nimbodex/microservices-factory/notification/internal/client/http/telegram"
	"github.com/nimbodex/microservices-factory/notification/internal/config"
	"github.com/nimbodex/microservices-factory/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/nimbodex/microservices-factory/notification/internal/service/consumer/order_paid_consumer"
	"github.com/nimbodex/microservices-factory/notification/internal/service/telegram"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/consumer"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
	"github.com/nimbodex/microservices-factory/platform/pkg/middleware/kafka"
)

func NewDIContainer(ctx context.Context, cfg *config.Config) (*DIContainer, error) {
	// Создаем логгер
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Создаем Telegram клиент
	telegramClient := telegram.NewClient(cfg.TelegramBot.GetBotToken())

	// Создаем Telegram сервис
	telegramService := telegram.NewService(telegramClient, log, cfg.TelegramBot.GetChatID())

	// Создаем Kafka consumer group для OrderPaid
	orderPaidConsumerConfig := sarama.NewConfig()
	orderPaidConsumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	orderPaidConsumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	orderPaidConsumerGroup, err := sarama.NewConsumerGroup(cfg.OrderPaidConsumer.GetBrokers(), cfg.OrderPaidConsumer.GetGroupID(), orderPaidConsumerConfig)
	if err != nil {
		return nil, err
	}

	// Создаем Kafka consumer group для OrderAssembled
	orderAssembledConsumerConfig := sarama.NewConfig()
	orderAssembledConsumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	orderAssembledConsumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	orderAssembledConsumerGroup, err := sarama.NewConsumerGroup(cfg.OrderAssembledConsumer.GetBrokers(), cfg.OrderAssembledConsumer.GetGroupID(), orderAssembledConsumerConfig)
	if err != nil {
		return nil, err
	}

	// Создаем Kafka consumers
	orderPaidKafkaConsumer := consumer.NewConsumer(
		orderPaidConsumerGroup,
		cfg.OrderPaidConsumer.GetTopics(),
		log,
		kafka.Logging(log),
	)

	orderAssembledKafkaConsumer := consumer.NewConsumer(
		orderAssembledConsumerGroup,
		cfg.OrderAssembledConsumer.GetTopics(),
		log,
		kafka.Logging(log),
	)

	// Создаем сервисы
	orderPaidConsumerService := order_paid_consumer.NewConsumerService(log, telegramService)
	orderAssembledConsumerService := order_assembled_consumer.NewConsumerService(log, telegramService)

	return &DIContainer{
		logger:                        log,
		telegramService:               telegramService,
		orderPaidKafkaConsumer:        orderPaidKafkaConsumer,
		orderAssembledKafkaConsumer:   orderAssembledKafkaConsumer,
		orderPaidConsumerGroup:        orderPaidConsumerGroup,
		orderAssembledConsumerGroup:   orderAssembledConsumerGroup,
		orderPaidConsumerService:      orderPaidConsumerService,
		orderAssembledConsumerService: orderAssembledConsumerService,
	}, nil
}

type DIContainer struct {
	logger                        logger.Logger
	telegramService               *telegram.Service
	orderPaidKafkaConsumer        consumer.Consumer
	orderAssembledKafkaConsumer   consumer.Consumer
	orderPaidConsumerGroup        sarama.ConsumerGroup
	orderAssembledConsumerGroup   sarama.ConsumerGroup
	orderPaidConsumerService      *order_paid_consumer.ConsumerService
	orderAssembledConsumerService *order_assembled_consumer.ConsumerService
}

func (c *DIContainer) GetLogger() logger.Logger {
	return c.logger
}

func (c *DIContainer) GetOrderPaidKafkaConsumer() consumer.Consumer {
	return c.orderPaidKafkaConsumer
}

func (c *DIContainer) GetOrderAssembledKafkaConsumer() consumer.Consumer {
	return c.orderAssembledKafkaConsumer
}

func (c *DIContainer) GetOrderPaidConsumerService() *order_paid_consumer.ConsumerService {
	return c.orderPaidConsumerService
}

func (c *DIContainer) GetOrderAssembledConsumerService() *order_assembled_consumer.ConsumerService {
	return c.orderAssembledConsumerService
}

func (c *DIContainer) Close() error {
	if err := c.orderPaidConsumerGroup.Close(); err != nil {
		return err
	}
	if err := c.orderAssembledConsumerGroup.Close(); err != nil {
		return err
	}
	return nil
}
