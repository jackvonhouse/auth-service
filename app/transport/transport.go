package transport

import (
	"github.com/gorilla/mux"
	"github.com/jackvonhouse/auth-service/app/usecase"
	"github.com/jackvonhouse/auth-service/internal/transport/auth"
	"github.com/jackvonhouse/auth-service/internal/transport/router"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type Transport struct {
	router *router.Router
}

func New(
	useCase *usecase.UseCase,
	logger log.Logger,
) *Transport {

	transportLogger := logger.WithField("layer", "transport")

	r := router.New("/api/v1")

	r.Handle(map[string]router.Handlify{
		"": auth.New(useCase.Auth, transportLogger),
	})

	return &Transport{
		router: r,
	}
}

func (t *Transport) Router() *mux.Router { return t.router.Router() }
