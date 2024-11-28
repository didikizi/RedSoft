package app

import (
	"context"
	"log/slog"

	"github.com/didikizi/RedSoft/iternal/config"
	"github.com/didikizi/RedSoft/iternal/router"
	"github.com/didikizi/RedSoft/iternal/service"
	"github.com/didikizi/RedSoft/iternal/storage"
)

func Start() error {
	ctx := context.Background()

	config, err := config.New()
	if err != nil {
		return err
	}

	slog.SetLogLoggerLevel(slog.Level(config.LogLevel))

	storage, err := storage.New(ctx, config)
	if err != nil {
		return err
	}

	service := service.New(storage)
	router := router.New(config, service)

	router.Start()

	return nil
}
