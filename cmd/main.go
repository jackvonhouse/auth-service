package main

import (
	"context"
	"flag"

	"github.com/jackvonhouse/auth-service/app"
	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/pkg/log"
	"github.com/jackvonhouse/auth-service/pkg/shutdown"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogrusLogger()

	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"config/config.toml",
		"The path to the configuration file",
	)

	flag.Parse()

	logger.Info("reading config...")

	config, err := config.New(configPath, logger)
	if err != nil {
		logger.Error(err)

		return
	}

	logger.Info("starting app...")

	app := app.New(ctx, config, logger)

	go app.Run()

	shutdown.Graceful(ctx, cancel, logger, app)
}
