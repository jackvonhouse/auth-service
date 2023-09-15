package app

import (
	"context"

	"github.com/jackvonhouse/auth-service/app/infrastructure"
	"github.com/jackvonhouse/auth-service/app/repository"
	"github.com/jackvonhouse/auth-service/app/service"
	"github.com/jackvonhouse/auth-service/app/transport"
	"github.com/jackvonhouse/auth-service/app/usecase"
	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/internal/infrastructure/server/http"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type App struct {
	infrastructure *infrastructure.Infrastructure
	repository     *repository.Repository
	service        *service.Service
	useCase        *usecase.UseCase
	transport      *transport.Transport

	config *config.Config
	logger log.Logger
	server *http.Server
}

func New(
	ctx context.Context,
	config *config.Config,
	logger log.Logger,
) *App {

	logger.Info("creating layers...")

	i := infrastructure.New(ctx, config, logger)
	r := repository.New(i, logger)
	s := service.New(r, config.JWT, logger)
	u := usecase.New(s, logger)
	t := transport.New(u, logger)

	logger.Info("creating http server...")

	httpServer := http.New(t.Router(), config.Server)

	return &App{
		infrastructure: i,
		repository:     r,
		service:        s,
		useCase:        u,
		transport:      t,
		config:         config,
		logger:         logger,
		server:         httpServer,
	}
}

func (a *App) Run() error {
	a.logger.Info("running http server...")

	return a.server.Run()
}

func (a *App) Shutdown(
	ctx context.Context,
) error {

	a.logger.Info("http server shutdowning..")

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.logger.Info("repository shutdowning..")

	if err := a.repository.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
