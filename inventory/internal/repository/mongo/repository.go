package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
)

// MongoPartRepository implements PartRepository using MongoDB
type MongoPartRepository struct {
	collection *mongo.Collection
}

// NewPartRepository creates a new MongoDB part repository
func NewPartRepository(client *mongo.Client, database string) *MongoPartRepository {
	return &MongoPartRepository{
		collection: client.Database(database).Collection("parts"),
	}
}

func (r *MongoPartRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	var part model.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid.String()}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("part with UUID %s not found", uuid)
		}
		return nil, fmt.Errorf("failed to get part: %w", err)
	}
	return &part, nil
}

func (r *MongoPartRepository) List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	var parts []*model.Part

	mongoFilter := bson.M{}
	if filter != nil {
		if len(filter.UUIDs) > 0 {
			mongoFilter["uuid"] = bson.M{"$in": filter.UUIDs}
		}
		if len(filter.Names) > 0 {
			mongoFilter["name"] = bson.M{"$in": filter.Names}
		}
		if len(filter.Categories) > 0 {
			mongoFilter["category"] = bson.M{"$in": filter.Categories}
		}
		if len(filter.ManufacturerCountries) > 0 {
			mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
		}
		if len(filter.Tags) > 0 {
			mongoFilter["tags"] = bson.M{"$in": filter.Tags}
		}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to find parts: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx) //nolint:gosec
	}()

	if err = cursor.All(ctx, &parts); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	return parts, nil
}

func (r *MongoPartRepository) Create(ctx context.Context, part *model.Part) error {
	if part == nil {
		return fmt.Errorf("part cannot be nil")
	}

	_, err := r.collection.InsertOne(ctx, part)
	if err != nil {
		return fmt.Errorf("failed to create part: %w", err)
	}

	return nil
}

func (r *MongoPartRepository) Update(ctx context.Context, part *model.Part) error {
	if part == nil {
		return fmt.Errorf("part cannot be nil")
	}

	filter := bson.M{"uuid": part.UUID}
	update := bson.M{"$set": part}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update part: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("part with UUID %s not found", part.UUID)
	}

	return nil
}

func (r *MongoPartRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	filter := bson.M{"uuid": uuid.String()}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete part: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("part with UUID %s not found", uuid)
	}

	return nil
}
