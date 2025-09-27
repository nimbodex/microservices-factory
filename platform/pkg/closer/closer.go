package closer

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
}

var (
	closerInstance *Closer
	once           sync.Once
)

type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(context.Context) error
	logger Logger
}

func New() *Closer {
	once.Do(func() {
		closerInstance = &Closer{
			done:  make(chan struct{}),
			funcs: make([]func(context.Context) error, 0),
		}
	})
	return closerInstance
}

func SetLogger(logger Logger) {
	closer := New()
	closer.mu.Lock()
	defer closer.mu.Unlock()
	closer.logger = logger
}

func AddNamed(name string, fn func(context.Context) error) {
	closer := New()
	closer.mu.Lock()
	defer closer.mu.Unlock()

	wrappedFn := func(ctx context.Context) error {
		if closer.logger != nil {
			closer.logger.Info(ctx, "Closing resource", "name", name)
		}
		if err := fn(ctx); err != nil {
			if closer.logger != nil {
				closer.logger.Error(ctx, "Failed to close resource", "name", name, "error", err)
			}
			return err
		}
		if closer.logger != nil {
			closer.logger.Info(ctx, "Resource closed successfully", "name", name)
		}
		return nil
	}

	closer.funcs = append(closer.funcs, wrappedFn)
}

func Add(fn func(context.Context) error) {
	AddNamed("unnamed", fn)
}

func CloseAll(ctx context.Context) error {
	closer := New()
	var err error

	closer.once.Do(func() {
		close(closer.done)

		for i := len(closer.funcs) - 1; i >= 0; i-- {
			if closeErr := closer.funcs[i](ctx); closeErr != nil {
				if err == nil {
					err = closeErr
				}
			}
		}
	})

	return err
}

func Configure(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	closer := New()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	go func() {
		sig := <-sigChan
		closer.logger.Info(context.Background(), "Received signal, shutting down", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()

		if err := CloseAll(ctx); err != nil {
			closer.logger.Error(ctx, "Error during shutdown", "error", err)
			os.Exit(1)
		}

		closer.logger.Info(ctx, "Shutdown completed")
		os.Exit(0)
	}()
}
