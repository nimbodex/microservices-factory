package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/nimbodex/microservices-factory/order/internal/config"
	"github.com/nimbodex/microservices-factory/platform/pkg/closer"
	"github.com/nimbodex/microservices-factory/platform/pkg/grpc/health"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

// App представляет основное приложение Order сервиса
type App struct {
	config     *config.Config
	logger     logger.Logger
	db         *sql.DB
	httpServer *http.Server
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

	if err := a.startServers(); err != nil {
		return fmt.Errorf("failed to start servers: %w", err)
	}

	a.logger.Info(ctx, "Order service started successfully")
	return nil
}

// initLogger инициализирует логгер
// loggerAdapter адаптирует logger.Logger для использования с closer
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

// initComponents инициализирует все компоненты приложения
func (a *App) initComponents(ctx context.Context) error {
	// Инициализируем PostgreSQL подключение
	if err := a.initPostgreSQL(ctx); err != nil {
		return fmt.Errorf("failed to init PostgreSQL: %w", err)
	}

	// Инициализируем gRPC сервер
	a.grpcServer = grpc.NewServer()
	health.RegisterService(a.grpcServer)

	// Инициализируем HTTP сервер
	a.httpServer = &http.Server{
		Addr:         a.config.OrderHTTP.Address(),
		ReadTimeout:  parseDuration(a.config.OrderHTTP.ReadTimeout()),
		WriteTimeout: 30 * time.Second,
		Handler:      a.createHTTPHandler(),
	}

	return nil
}

// initPostgreSQL инициализирует подключение к PostgreSQL
func (a *App) initPostgreSQL(ctx context.Context) error {
	db, err := sql.Open("postgres", a.config.Postgres.URI())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Проверяем подключение
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	a.db = db
	a.logger.Info(ctx, "PostgreSQL connected successfully")

	return nil
}

func (a *App) startServers() error {
	go func() {
		a.logger.Info(context.Background(), "HTTP server starting",
			zap.String("address", a.config.OrderHTTP.Address()))

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error(context.Background(), "HTTP server failed", zap.Error(err))
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":50052") // Другой порт для Order gRPC
		if err != nil {
			a.logger.Error(context.Background(), "Failed to listen for gRPC", zap.Error(err))
			return
		}

		a.logger.Info(context.Background(), "gRPC server starting",
			zap.String("address", ":50052"))

		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error(context.Background(), "gRPC server failed", zap.Error(err))
		}
	}()

	return nil
}

func (a *App) createHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}

// setupGracefulShutdown настраивает graceful shutdown
func (a *App) setupGracefulShutdown() {
	// Регистрируем функции закрытия
	closer.AddNamed("HTTP Server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	closer.AddNamed("gRPC Server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	closer.AddNamed("PostgreSQL", func(ctx context.Context) error {
		if a.db != nil {
			return a.db.Close()
		}
		return nil
	})

	// Настраиваем обработку сигналов
	closer.Configure(syscall.SIGTERM, syscall.SIGINT)

}

// parseDuration парсит строку в time.Duration
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * time.Second // значение по умолчанию
	}
	return d
}
