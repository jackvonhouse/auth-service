package infrastructure

import (
	"context"

	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/internal/infrastructure/mongo"
	"github.com/jackvonhouse/auth-service/pkg/log"
)

type Infrastructure struct {
	Mongo *mongo.DatabaseMongo
}

func New(
	ctx context.Context,
	config *config.Config,
	logger log.Logger,
) *Infrastructure {

	infrastructureLog := logger.WithField("layer", "infrastructure")

	infrastructureLog.Info("opening mongo connection...")

	db, err := mongo.New(ctx, config.Database, infrastructureLog)
	if err != nil {
		infrastructureLog.Warn(err)
	}

	return &Infrastructure{
		Mongo: db,
	}
}
