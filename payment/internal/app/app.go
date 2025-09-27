package app

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/nimbodex/microservices-factory/payment/internal/config"
	"github.com/nimbodex/microservices-factory/platform/pkg/closer"
	"github.com/nimbodex/microservices-factory/platform/pkg/grpc/health"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

type App struct {
	config     *config.Config
	logger     logger.Logger
	grpcServer *grpc.Server
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

	if err := a.initComponents(ctx); err != nil {
		return fmt.Errorf("failed to init components: %w", err)
	}

	a.setupGracefulShutdown()

	if err := a.startGRPCServer(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	a.logger.Info(ctx, "Payment service started successfully")
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
	err := logger.Init(a.config.Logger.Level(), a.config.Logger.AsJSON())
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	a.logger = logger.GetLogger()
	closer.SetLogger(&loggerAdapter{logger: a.logger})
	return nil
}

func (a *App) initComponents(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()

	health.RegisterService(a.grpcServer)

	return nil
}

func (a *App) startGRPCServer() error {
	lis, err := net.Listen("tcp", a.config.PaymentGRPC.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", a.config.PaymentGRPC.Address(), err)
	}

	go func() {
		a.logger.Info(context.Background(), "gRPC server starting",
			zap.String("address", a.config.PaymentGRPC.Address()))

		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error(context.Background(), "gRPC server failed", zap.Error(err))
		}
	}()

	return nil
}

func (a *App) setupGracefulShutdown() {
	closer.AddNamed("gRPC Server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	closer.Configure(syscall.SIGTERM, syscall.SIGINT)
}
