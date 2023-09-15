package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackvonhouse/auth-service/internal/dto"
	"github.com/jackvonhouse/auth-service/internal/transport"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type useCaseAuth interface {
	Auth(context.Context, *dto.Auth) (*dto.TokenPair, error)
	Refresh(context.Context, *dto.TokenPair) (*dto.TokenPair, error)
}

type TransportAuth struct {
	logger  log.Logger
	useCase useCaseAuth
}

func New(
	useCase useCaseAuth,
	logger log.Logger,
) *TransportAuth {

	return &TransportAuth{
		useCase: useCase,
		logger:  logger.WithField("layer", "transport"),
	}
}

func (t *TransportAuth) Handle(
	router *mux.Router,
) {

	router.HandleFunc("/auth", t.Auth).
		Methods(http.MethodPost).
		Queries("guid", "{guid}")

	router.HandleFunc("/refresh", t.Refresh).
		Methods(http.MethodPost)
}

func (t *TransportAuth) Auth(
	w http.ResponseWriter,
	r *http.Request,
) {

	queries := r.URL.Query()

	guid := queries.Get("guid")
	if guid == "" {
		t.logger.Warn("guid is empty")

		transport.Error(w, http.StatusBadRequest, "guid is empty")

		return
	}

	// if !validator.IsValidGUID(guid) {
	// 	t.logger.WithField("guid", guid).
	// 		Warn("invalid guid")

	// 	transport.Error(w, http.StatusBadRequest, "invalid guid")

	// 	return
	// }

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data := dto.Auth{
		GUID: guid,
	}

	tokenPair, err := t.useCase.Auth(ctx, &data)
	if err != nil {
		t.logger.Warn(err)

		transport.Error(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)

		return
	}

	transport.Response(w, tokenPair)
}

func (t *TransportAuth) Refresh(
	w http.ResponseWriter,
	r *http.Request,
) {

	data := dto.TokenPair{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		t.logger.Warn(err)

		transport.Error(w, http.StatusBadRequest, "invalid json structure")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tokenPair, err := t.useCase.Refresh(ctx, &data)
	if err != nil {
		t.logger.Warn(err)

		code, msg := 0, ""

		if strings.Contains(err.Error(), "ErrInternal") {
			code = http.StatusInternalServerError
			msg = http.StatusText(code)
		} else {
			code = http.StatusUnauthorized
			msg = err.Error()
		}

		transport.Error(w, code, msg)

		return
	}

	transport.Response(w, tokenPair)
}
