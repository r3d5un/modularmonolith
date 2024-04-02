package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/r3d5un/modularmonolith/cmd/monolith/rabbit"
	"github.com/r3d5un/modularmonolith/cmd/monolith/warehouse"
	"github.com/r3d5un/modularmonolith/internal/config"
	"github.com/r3d5un/modularmonolith/internal/monolith"
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

	logger.Info("constructing application")
	app := application{logger: logger}

	logger.Info("loading configuration")
	app.cfg, err = config.New()
	if err != nil {
		bareLogger.Error("unable to load config", "error", err)
		return err
	}
	logger.Info("configuration loaded", "configuration", app.cfg)

	logger.Info("opening database connection pool...")
	app.db, err = openDB(app.cfg.DB)
	if err != nil {
		logger.Error("error occurred while connecting to the database",
			"error", err,
		)
		return err
	}
	defer app.db.Close()
	logger.Info("database connection pool established")

	logger.Info("opening message queue connection pool...")
	app.mq, err = openQueue(app.cfg.MQ)
	if err != nil {
		slog.Error("unable to create message queue connection pool", "error", err)
		return err
	}
	defer app.mq.Shutdown()
	logger.Info("message queue connection pool established")

	app.mux = http.NewServeMux()

	logger.Info("initializing modules")
	app.modules = &monolith.Modules{
		Warehouse: &warehouse.Module{},
		Rabbit:    &rabbit.Module{},
	}

	logger.Info("running module startup prodedures")
	app.setupModules(context.Background())
	app.postSetupModules(context.Background())

	err = app.serve()
	if err != nil {
		logger.Error("unable to start server", "error", err)
		return err
	}

	return nil
}

func (app *application) setupModules(ctx context.Context) {
	val := reflect.ValueOf(app.modules)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if module, ok := field.Interface().(monolith.Module); ok {
			module.Setup(ctx, app)
		}
	}
}

func (app *application) postSetupModules(ctx context.Context) {
	val := reflect.ValueOf(app.modules)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if module, ok := field.Interface().(monolith.Module); ok {
			module.PostSetup(ctx, app)
		}
	}
}
