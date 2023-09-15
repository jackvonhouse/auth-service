package mongo

import (
	"context"
	"fmt"

	"github.com/jackvonhouse/auth-service/config"
	"github.com/jackvonhouse/auth-service/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseMongo struct {
	client *mongo.Client
}

func New(
	ctx context.Context,
	config *config.Database,
	logger log.Logger,
) (*DatabaseMongo, error) {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	client, err := mongo.Connect(
		ctx,

		options.Client().
			ApplyURI(config.String()).
			SetServerAPIOptions(serverAPI),
	)

	if err != nil {
		logger.Warnf("can't connect to mongo: %s", err)

		return nil, fmt.Errorf("can't connect to mongo: %s", err)
	}

	return &DatabaseMongo{
		client: client,
	}, nil
}

func (d *DatabaseMongo) Database() *mongo.Client { return d.client }
