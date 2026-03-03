package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wbtech/internal/config"
	"wbtech/internal/infrastructure/cache"
	httpServer "wbtech/internal/infrastructure/http"
	"wbtech/internal/infrastructure/http/handler"
	"wbtech/internal/infrastructure/kafka"
	"wbtech/internal/infrastructure/postgres"
	"wbtech/internal/usecase/order"
	_ "wbtech/metrics"
)

type App struct {
	httpSrv  *http.Server
	useCase  *order.OrderUseCase
	consumer *kafka.Consumer
	db       *sql.DB
}

func NewApp(cfg config.Config) (*App, error) {
	db, err := postgres.NewPostgresConnection(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		return nil, err
	}
	repo := postgres.NewOrderRepository(db)

	orderCache := cache.NewOrderCache(100)
	useCase := order.NewOrderUseCase(repo, orderCache)
	if err := useCase.RestoreCache(context.Background()); err != nil {
		log.Printf("failed to restore cache: %v", err)
	}
	orderHandler := handler.NewOrderHandler(useCase)
	router := httpServer.NewRouter(orderHandler)

	httpSrv := httpServer.NewServer(":"+cfg.HTTPPort, router)

	consumer := kafka.NewConsumer(
		cfg.KafkaBroker,
		cfg.KafkaTopic,
		cfg.KafkaGroup,
		useCase,
	)

	return &App{
		httpSrv:  httpSrv,
		useCase:  useCase,
		consumer: consumer,
		db:       db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go a.consumer.Start(ctx)
	go func() {
		log.Printf("HTTP server started on %s", a.httpSrv.Addr)
		if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("shutting down by context")
	case err := <-errCh:
		log.Printf("HTTP server error: %v", err)
		return err
	case sig := <-sigCh:
		log.Printf("received signal: %v", sig)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if a.consumer != nil {
		_ = a.consumer.Close()
	}
	if err := a.db.Close(); err != nil {
		return err
	}
	log.Println("server stopped gracefully")
	return nil
}
