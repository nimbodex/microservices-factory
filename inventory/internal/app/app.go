package app

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	v1 "github.com/nimbodex/microservices-factory/inventory/internal/api/inventory/v1"
	"github.com/nimbodex/microservices-factory/inventory/internal/config"
	"github.com/nimbodex/microservices-factory/inventory/internal/repository"
	"github.com/nimbodex/microservices-factory/inventory/internal/service"
	"github.com/nimbodex/microservices-factory/platform/pkg/closer"
	"github.com/nimbodex/microservices-factory/platform/pkg/grpc/health"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

type App struct {
	config         *config.Config
	logger         logger.Logger
	mongoClient    *mongo.Client
	grpcServer     *grpc.Server
	partRepository repository.PartRepository
	partService    service.InventoryService
	grpcService    *v1.APIHandler
}

func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run() error {
	ctx := context.Background()

	if err := a.initLogger(); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	a.setupGracefulShutdown()

	if err := a.initServices(ctx); err != nil {
		return fmt.Errorf("failed to init services: %w", err)
	}

	if err := a.startGRPCServer(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	a.logger.Info(ctx, "Inventory service started successfully")
	return nil
}

type loggerAdapter struct {
	logger logger.Logger
}

func (l *loggerAdapter) Info(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any("field", field)
	}
	l.logger.Info(ctx, msg, zapFields...)
}

func (l *loggerAdapter) Error(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any("field", field)
	}
	l.logger.Error(ctx, msg, zapFields...)
}

func (a *App) initLogger() error {
	if err := logger.Init(a.config.Logger.Level(), a.config.Logger.AsJSON()); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	a.logger = logger.GetLogger()
	closer.SetLogger(&loggerAdapter{logger: a.logger})
	return nil
}

func (a *App) setupGracefulShutdown() {
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	closer.Add(func(ctx context.Context) error {
		if a.grpcServer != nil {
			a.grpcServer.GracefulStop()
			a.logger.Info(ctx, "gRPC server stopped gracefully")
		}
		return nil
	})
	closer.Add(func(ctx context.Context) error {
		if a.mongoClient != nil {
			if err := a.mongoClient.Disconnect(ctx); err != nil {
				return fmt.Errorf("failed to disconnect MongoDB: %w", err)
			}
			a.logger.Info(ctx, "MongoDB disconnected")
		}
		return nil
	})
}

func (a *App) initServices(ctx context.Context) error {
	if err := a.initMongoDB(ctx); err != nil {
		return fmt.Errorf("failed to init MongoDB: %w", err)
	}

	a.partRepository = repository.NewPartRepository(a.mongoClient, a.config.Mongo.Database())
	a.partService = service.NewInventoryService(a.partRepository)
	a.grpcService = v1.NewAPIHandler(a.partService)
	a.grpcServer = grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(a.grpcServer, a.grpcService)
	health.RegisterService(a.grpcServer)

	return nil
}

func (a *App) initMongoDB(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(a.config.Mongo.URI())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	a.mongoClient = client
	a.logger.Info(ctx, "MongoDB connected successfully")
	return nil
}

func (a *App) startGRPCServer() error {
	lis, err := net.Listen("tcp", a.config.InventoryGRPC.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", a.config.InventoryGRPC.Address(), err)
	}

	go func() {
		a.logger.Info(context.Background(), "gRPC server starting",
			zap.String("address", a.config.InventoryGRPC.Address()))

		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error(context.Background(), "gRPC server failed", zap.Error(err))
		}
	}()
	return nil
}
