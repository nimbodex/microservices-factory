package part

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// MongoPartRepository implements PartRepository using MongoDB
type MongoPartRepository struct {
	collection *mongo.Collection
}

// NewMongoPartRepository creates a new MongoDB part repository
func NewMongoPartRepository(collection *mongo.Collection) *MongoPartRepository {
	return &MongoPartRepository{
		collection: collection,
	}
}

// GetByUUID retrieves a part by its UUID
func (r *MongoPartRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	filter := bson.M{"uuid": uuid.String()}

	var doc bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("part with UUID %s not found", uuid)
		}
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	part, err := r.documentToPart(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert document to part: %w", err)
	}

	return part, nil
}

// List retrieves parts matching the filter criteria
func (r *MongoPartRepository) List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	mongoFilter := r.buildMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to find parts: %w", err)
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			log.Printf("failed to close cursor: %s", closeErr.Error())
		}
	}()

	var parts []*model.Part
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		part, err := r.documentToPart(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert document to part: %w", err)
		}

		parts = append(parts, part)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return parts, nil
}

// Create creates a new part
func (r *MongoPartRepository) Create(ctx context.Context, part *model.Part) error {
	doc := r.partToDocument(part)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("part with UUID %s already exists", part.UUID)
		}
		return fmt.Errorf("failed to create part: %w", err)
	}

	return nil
}

// Update updates an existing part
func (r *MongoPartRepository) Update(ctx context.Context, part *model.Part) error {
	filter := bson.M{"uuid": part.UUID.String()}
	update := bson.M{"$set": r.partToDocument(part)}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update part: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("part with UUID %s not found", part.UUID)
	}

	return nil
}

// Delete removes a part by its UUID
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

// buildMongoFilter builds MongoDB filter from service filter
func (r *MongoPartRepository) buildMongoFilter(filter *model.PartsFilter) bson.M {
	if isEmptyFilter(filter) {
		return bson.M{}
	}

	mongoFilter := bson.M{}

	if len(filter.UUIDs) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		nameRegexes := make([]bson.M, len(filter.Names))
		for i, name := range filter.Names {
			nameRegexes[i] = bson.M{
				"name": bson.M{
					"$regex":   primitive.Regex{Pattern: strings.ToLower(name), Options: "i"},
					"$options": "i",
				},
			}
		}
		mongoFilter["$or"] = nameRegexes
	}

	if len(filter.Categories) > 0 {
		categories := make([]int32, len(filter.Categories))
		for i, category := range filter.Categories {
			categories[i] = int32(category)
		}
		mongoFilter["category"] = bson.M{"$in": categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{
			"$in": filter.ManufacturerCountries,
		}
	}

	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{
			"$in": filter.Tags,
		}
	}

	return mongoFilter
}

// partToDocument converts Part model to MongoDB document
func (r *MongoPartRepository) partToDocument(part *model.Part) bson.M {
	doc := bson.M{
		"uuid":           part.UUID.String(),
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       int32(part.Category),
		"tags":           part.Tags,
		"metadata":       part.Metadata,
		"created_at":     part.CreatedAt,
		"updated_at":     part.UpdatedAt,
	}

	if part.Dimensions != nil {
		doc["dimensions"] = bson.M{
			"length": part.Dimensions.Length,
			"width":  part.Dimensions.Width,
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		}
	}

	if part.Manufacturer != nil {
		doc["manufacturer"] = bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		}
	}

	return doc
}

// documentToPart converts MongoDB document to Part model
func (r *MongoPartRepository) documentToPart(doc bson.M) (*model.Part, error) { //nolint
	part := &model.Part{}

	if uuidStr, ok := doc["uuid"].(string); ok {
		parsedUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID format: %w", err)
		}
		part.UUID = parsedUUID
	}

	if name, ok := doc["name"].(string); ok {
		part.Name = name
	}
	if description, ok := doc["description"].(string); ok {
		part.Description = description
	}
	if price, ok := doc["price"].(float64); ok {
		part.Price = price
	}
	if stockQuantity, ok := doc["stock_quantity"].(int32); ok {
		part.StockQuantity = stockQuantity
	}
	if category, ok := doc["category"].(int32); ok {
		part.Category = inventoryv1.Category(category)
	}

	if createdAt, ok := doc["created_at"].(primitive.DateTime); ok {
		part.CreatedAt = createdAt.Time()
	}
	if updatedAt, ok := doc["updated_at"].(primitive.DateTime); ok {
		part.UpdatedAt = updatedAt.Time()
	}

	if tags, ok := doc["tags"].(primitive.A); ok {
		part.Tags = make([]string, len(tags))
		for i, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				part.Tags[i] = tagStr
			}
		}
	}

	if metadata, ok := doc["metadata"].(bson.M); ok {
		part.Metadata = make(map[string]interface{})
		for k, v := range metadata {
			part.Metadata[k] = v
		}
	}

	if dimensionsDoc, ok := doc["dimensions"].(bson.M); ok {
		part.Dimensions = &model.Dimensions{}
		if length, ok := dimensionsDoc["length"].(float64); ok {
			part.Dimensions.Length = length
		}
		if width, ok := dimensionsDoc["width"].(float64); ok {
			part.Dimensions.Width = width
		}
		if height, ok := dimensionsDoc["height"].(float64); ok {
			part.Dimensions.Height = height
		}
		if weight, ok := dimensionsDoc["weight"].(float64); ok {
			part.Dimensions.Weight = weight
		}
	}

	if manufacturerDoc, ok := doc["manufacturer"].(bson.M); ok {
		part.Manufacturer = &model.Manufacturer{}
		if name, ok := manufacturerDoc["name"].(string); ok {
			part.Manufacturer.Name = name
		}
		if country, ok := manufacturerDoc["country"].(string); ok {
			part.Manufacturer.Country = country
		}
		if website, ok := manufacturerDoc["website"].(string); ok {
			part.Manufacturer.Website = website
		}
	}

	return part, nil
}
