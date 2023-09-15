package repository

import (
	"context"

	"github.com/jackvonhouse/auth-service/app/infrastructure"
	"github.com/jackvonhouse/auth-service/internal/infrastructure/mongo"
	"github.com/jackvonhouse/auth-service/internal/repository/auth"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type Repository struct {
	Auth  *auth.RepositoryAuth
	Mongo *mongo.DatabaseMongo
}

func New(
	infrastructure *infrastructure.Infrastructure,
	logger log.Logger,
) *Repository {

	repositoryLogger := logger.WithField("layer", "repository")

	return &Repository{
		Auth: auth.New(
			infrastructure.Mongo.Database(),
			repositoryLogger,
		),
		Mongo: infrastructure.Mongo,
	}
}

func (r *Repository) Shutdown(
	ctx context.Context,
) error {

	return r.Mongo.Database().Disconnect(ctx)
}
