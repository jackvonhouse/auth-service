package auth

import (
	"context"
	"fmt"

	"github.com/jackvonhouse/auth-service/internal/dto"
	"github.com/jackvonhouse/auth-service/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoryAuth struct {
	database *mongo.Client
	logger   log.Logger
}

func New(
	client *mongo.Client,
	logger log.Logger,
) *RepositoryAuth {

	client.Database("auth").
		Collection("refresh_tokens").
		Indexes().
		CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{
				"expire_at": 1,
			},
			Options: options.Index().SetExpireAfterSeconds(0),
		})

	return &RepositoryAuth{
		database: client,
		logger:   logger.WithField("unit", "auth"),
	}
}

func (r *RepositoryAuth) CreateRefreshToken(
	ctx context.Context,
	data *dto.RefreshToken,
) (string, error) {

	collection := r.database.Database("auth").
		Collection("refresh_tokens")

	result, err := collection.InsertOne(ctx, &data)
	if err != nil {
		r.logger.Warnf("can't insert refresh token: %s", err)

		return "", fmt.Errorf("ErrInternal: can't insert refresh token")
	}

	return fmt.Sprintf("%s",
		result.InsertedID.(primitive.ObjectID).Hex(),
	), nil
}

func (r *RepositoryAuth) GetRefreshToken(
	ctx context.Context,
	id string,
) (*dto.RefreshToken, error) {

	collection := r.database.Database("auth").
		Collection("refresh_tokens")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"id": id,
		}).Warnf("cant get object id from passed id: %s", err)

		return nil, fmt.Errorf("ErrInternal: cant get object id from passed id: %s", err)
	}

	filter := bson.D{{
		Key:   "_id",
		Value: objectId,
	}}

	data := dto.RefreshToken{}

	if err := collection.FindOne(ctx, filter).Decode(&data); err != nil {
		r.logger.WithFields(map[string]any{
			"id": id,
		}).Warnf("cant find refresh token by id: %s", err)

		return nil, fmt.Errorf("refresh token not exists")
	}

	return &data, nil
}

func (r *RepositoryAuth) DeleteRefreshToken(
	ctx context.Context,
	id string,
) error {

	collection := r.database.Database("auth").
		Collection("refresh_tokens")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"id": id,
		}).Warnf("cant get object id from passed id: %s", err)

		return fmt.Errorf("ErrInternal: cant get object id from passed id")
	}

	filter := bson.D{{
		Key:   "_id",
		Value: objectId,
	}}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"id": id,
		}).Warnf("cant delete refresh token by id: %s", err)

		return fmt.Errorf("ErrInternal: cant delete refresh token by id")
	}

	return nil
}
