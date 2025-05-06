package main

import (
	"context"
	"fmt"
	defaultLogger "log"
	"os"
	"path/filepath"
	"wallets-service/internal/wallets/connections"

	kafkaPkg "github.com/knstch/subtrack-kafka/consumer"
	"github.com/knstch/subtrack-libs/log"
	"github.com/knstch/subtrack-libs/tracing"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wallets-service/internal/endpoints/consumer"
	"wallets-service/internal/wallets"
	"wallets-service/internal/wallets/repo"

	"wallets-service/config"
)

func main() {
	if err := run(); err != nil {
		defaultLogger.Print(err)
		recover()
	}
}

func run() error {
	args := os.Args

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	if err := config.InitENV(dir); err != nil {
		return fmt.Errorf("config.InitENV: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig: %w", err)
	}

	shutdown := tracing.InitTracer(cfg.ServiceName, cfg.JaegerHost)
	defer shutdown(context.Background())

	logger := log.NewLogger(cfg.ServiceName, log.InfoLevel)

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}
	dbRepo := repo.NewDBRepo(logger, db)

	conns, err := connections.MakeConnections(*cfg, logger)
	if err != nil {
		return fmt.Errorf("connections.MakeConnections: %w", err)
	}

	svc := wallets.NewService(logger, dbRepo, cfg, conns, nil)

	walletsConsumer, err := kafkaPkg.NewConsumer(cfg.KafkaAddr, "wallets-group", logger)
	if err != nil {
		return fmt.Errorf("consumer.NewConsumer: %w", err)
	}

	consumerController := consumer.NewController(logger, svc)

	consumerController.InitHandlers(walletsConsumer)

	if err = walletsConsumer.Run(context.Background()); err != nil {
		return fmt.Errorf("walletsConsumer.Run: %w", err)
	}

	return nil
}
