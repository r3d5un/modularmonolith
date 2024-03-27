package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/r3d5un/modularmonolith/cmd/monolith/config"
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

	return nil
}
