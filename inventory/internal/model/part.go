package model

import (
	"time"

	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// Part represents a part in the service layer
type Part struct {
	UUID          string                 `json:"uuid" bson:"uuid"`
	Name          string                 `json:"name" bson:"name"`
	Description   string                 `json:"description" bson:"description"`
	Price         float64                `json:"price" bson:"price"`
	StockQuantity int32                  `json:"stock_quantity" bson:"stock_quantity"`
	Category      inventoryv1.Category   `json:"category" bson:"category"`
	Dimensions    *Dimensions            `json:"dimensions" bson:"dimensions"`
	Manufacturer  *Manufacturer          `json:"manufacturer" bson:"manufacturer"`
	Tags          []string               `json:"tags" bson:"tags"`
	Metadata      map[string]interface{} `json:"metadata" bson:"metadata"`
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" bson:"updated_at"`
}

// Dimensions represents part dimensions
type Dimensions struct {
	Length float64 `json:"length" bson:"length"`
	Width  float64 `json:"width" bson:"width"`
	Height float64 `json:"height" bson:"height"`
	Weight float64 `json:"weight" bson:"weight"`
}

// Manufacturer represents part manufacturer
type Manufacturer struct {
	Name    string `json:"name" bson:"name"`
	Country string `json:"country" bson:"country"`
	Website string `json:"website" bson:"website"`
}

// PartsFilter represents filter criteria for listing parts
type PartsFilter struct {
	UUIDs                 []string               `json:"uuids"`
	Names                 []string               `json:"names"`
	Categories            []inventoryv1.Category `json:"categories"`
	ManufacturerCountries []string               `json:"manufacturer_countries"`
	Tags                  []string               `json:"tags"`
}
