package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/maxim/microservices-factory/payment/internal/config"
	"github.com/maxim/microservices-factory/platform/pkg/closer"
	"github.com/maxim/microservices-factory/platform/pkg/grpc/health"
	"github.com/maxim/microservices-factory/platform/pkg/logger"
)

// App представляет основное приложение Payment сервиса
type App struct {
	config     *config.Config
	logger     *logger.Logger
	grpcServer *grpc.Server
}

// New создает новое приложение с заданной конфигурацией
func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run() error {
	ctx := context.Background()

	// Инициализируем логгер
	if err := a.initLogger(); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	// Инициализируем компоненты
	if err := a.initComponents(ctx); err != nil {
		return fmt.Errorf("failed to init components: %w", err)
	}

	// Настраиваем graceful shutdown
	a.setupGracefulShutdown()

	if err := a.startGRPCServer(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	a.logger.Info(ctx, "Payment service started successfully")
	return nil
}

// initLogger инициализирует логгер
func (a *App) initLogger() error {
	err := logger.Init(a.config.Logger.Level(), a.config.Logger.AsJSON())
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	a.logger = logger.Logger()
	return nil
}

// initComponents инициализирует все компоненты приложения
func (a *App) initComponents(ctx context.Context) error {
	// Инициализируем gRPC сервер
	a.grpcServer = grpc.NewServer()

	// Регистрируем сервисы
	// Здесь должен быть регистрирован Payment сервис
	// paymentv1.RegisterPaymentServiceServer(a.grpcServer, paymentService)
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

// setupGracefulShutdown настраивает graceful shutdown
func (a *App) setupGracefulShutdown() {
	// Регистрируем функции закрытия
	closer.AddNamed("gRPC Server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	// Настраиваем обработку сигналов
	closer.Configure(syscall.SIGTERM, syscall.SIGINT)

	// Обработка сигнала завершения
	go func() {
		<-closer.Done()
		a.logger.Info(context.Background(), "Application shutdown completed")
		os.Exit(0)
	}()
}
