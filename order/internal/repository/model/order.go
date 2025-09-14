package model

import (
	"time"
)

// Order represents an order in the repository layer
type Order struct {
	UUID      string    `json:"uuid"`
	UserUUID  string    `json:"user_uuid"`
	PartUUIDs []string  `json:"part_uuids"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Part represents a part in the repository layer
type Part struct {
	UUID  string  `json:"uuid"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
