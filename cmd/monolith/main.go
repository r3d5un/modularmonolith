package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/r3d5un/modularmonolith/internal/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	bareLogger := slog.New(handler)
	logger := bareLogger.With(
		slog.Group(
			"application_instance",
			slog.String("instance_id", uuid.New().String()),
		),
	)
	slog.SetDefault(logger)

	logger.Info("loading configuration")
	config, err := config.New()
	if err != nil {
		bareLogger.Error("unable to load config", "error", err)
		return err
	}
	logger.Info("configuration loaded", "configuration", config)

	logger.Info("opening database connection pool...")
	db, err := openDB(config.DB)
	if err != nil {
		logger.Error("error occurred while connecting to the database",
			"error", err,
		)
		return err
	}
	defer db.Close()
	logger.Info("database connection pool established")

	logger.Info("opening message queue connection pool...")
	mq, err := openQueue(config.MQ)
	if err != nil {
		slog.Error("unable to create message queue connection pool", "error", err)
		return err
	}
	defer mq.Shutdown()
	logger.Info("message queue connection pool established")

	app := application{
		cfg:    config,
		db:     db,
		mq:     mq,
		mux:    http.NewServeMux(),
		logger: logger,
	}

	logger.Info("starting server", "settings", config.App)
	err = app.serve()
	if err != nil {
		logger.Error("unable to start server", "error", err)
		return err
	}

	return nil
}
