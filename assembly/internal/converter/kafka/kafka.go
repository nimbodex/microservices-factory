package kafka

import (
	"encoding/json"

	"github.com/nimbodex/microservices-factory/assembly/internal/model"
)

// EncodeShipAssembledEvent кодирует ShipAssembledEvent в JSON
func EncodeShipAssembledEvent(event *model.ShipAssembledEvent) ([]byte, error) {
	return json.Marshal(event)
}
