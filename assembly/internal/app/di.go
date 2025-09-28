package app

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/nimbodex/microservices-factory/assembly/internal/config"
	"github.com/nimbodex/microservices-factory/assembly/internal/service"
	"github.com/nimbodex/microservices-factory/assembly/internal/service/consumer/order_consumer"
	"github.com/nimbodex/microservices-factory/assembly/internal/service/producer/order_producer"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/consumer"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/producer"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
	"github.com/nimbodex/microservices-factory/platform/pkg/middleware/kafka"
)

func NewDIContainer(ctx context.Context, cfg *config.Config) (*DIContainer, error) {
	// Создаем логгер
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Создаем Kafka consumer group
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.OrderPaidConsumer.GetBrokers(), cfg.OrderPaidConsumer.GetGroupID(), consumerConfig)
	if err != nil {
		return nil, err
	}

	// Создаем Kafka producer
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 3
	producerConfig.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(cfg.OrderAssembledProducer.GetBrokers(), producerConfig)
	if err != nil {
		return nil, err
	}

	// Создаем Kafka consumer
	kafkaConsumer := consumer.NewConsumer(
		consumerGroup,
		cfg.OrderPaidConsumer.GetTopics(),
		log,
		kafka.Logging(log),
	)

	// Создаем Kafka producer
	kafkaProducer := producer.NewProducer(
		syncProducer,
		cfg.OrderAssembledProducer.GetTopic(),
		log,
	)

	// Создаем сервисы
	orderProducerService := order_producer.NewProducerService(kafkaProducer, log)
	orderConsumerService := order_consumer.NewConsumerService(log, orderProducerService)

	// Создаем основной сервис
	service := service.NewService(orderConsumerService)

	return &DIContainer{
		logger:        log,
		kafkaConsumer: kafkaConsumer,
		kafkaProducer: kafkaProducer,
		consumerGroup: consumerGroup,
		syncProducer:  syncProducer,
		service:       service,
		orderProducer: orderProducerService,
		orderConsumer: orderConsumerService,
	}, nil
}

type DIContainer struct {
	logger        logger.Logger
	kafkaConsumer consumer.Consumer
	kafkaProducer producer.Producer
	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer
	service       *service.Service
	orderProducer *order_producer.ProducerService
	orderConsumer *order_consumer.ConsumerService
}

func (c *DIContainer) GetLogger() logger.Logger {
	return c.logger
}

func (c *DIContainer) GetKafkaConsumer() consumer.Consumer {
	return c.kafkaConsumer
}

func (c *DIContainer) GetService() *service.Service {
	return c.service
}

func (c *DIContainer) Close() error {
	if err := c.consumerGroup.Close(); err != nil {
		return err
	}
	if err := c.syncProducer.Close(); err != nil {
		return err
	}
	return nil
}
