package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/nimbodex/microservices-factory/order/internal/config"
	"github.com/nimbodex/microservices-factory/platform/pkg/closer"
	"github.com/nimbodex/microservices-factory/platform/pkg/grpc/health" //nolint
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

type App struct {
	config     *config.Config
	logger     logger.Logger
	db         *sql.DB
	httpServer *http.Server
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

	a.startServers()

	a.logger.Info(ctx, "Order service started successfully")
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
	if err := a.initPostgreSQL(ctx); err != nil {
		return fmt.Errorf("failed to init PostgreSQL: %w", err)
	}

	a.grpcServer = grpc.NewServer()
	health.RegisterService(a.grpcServer)

	a.httpServer = &http.Server{
		Addr:         a.config.OrderHTTP.Address(),
		ReadTimeout:  parseDuration(a.config.OrderHTTP.ReadTimeout()),
		WriteTimeout: 30 * time.Second,
		Handler:      a.createHTTPHandler(),
	}

	return nil
}

func (a *App) initPostgreSQL(ctx context.Context) error {
	db, err := sql.Open("postgres", a.config.Postgres.URI())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	a.db = db
	a.logger.Info(ctx, "PostgreSQL connected successfully")

	return nil
}

func (a *App) startServers() {
	go func() {
		a.logger.Info(context.Background(), "HTTP server starting",
			zap.String("address", a.config.OrderHTTP.Address()))

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error(context.Background(), "HTTP server failed", zap.Error(err))
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", "localhost:50052")
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
}

func (a *App) createHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK")) //nolint:gosec
	})

	return mux
}

func (a *App) setupGracefulShutdown() {
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

	closer.Configure(syscall.SIGTERM, syscall.SIGINT)
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * time.Second
	}
	return d
}
