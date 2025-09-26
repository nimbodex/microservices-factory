package model

import (
	"time"

	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// Part represents a part in the repository layer
type Part struct {
	UUID          string                 `json:"uuid"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Price         float64                `json:"price"`
	StockQuantity int32                  `json:"stock_quantity"`
	Category      inventoryv1.Category   `json:"category"`
	Dimensions    *Dimensions            `json:"dimensions"`
	Manufacturer  *Manufacturer          `json:"manufacturer"`
	Tags          []string               `json:"tags"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// Dimensions represents part dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
}

// Manufacturer represents part manufacturer
type Manufacturer struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
}
