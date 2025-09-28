package decoder

import (
	"encoding/json"

	"github.com/nimbodex/microservices-factory/order/internal/model"
)

// DecodeShipAssembledEvent декодирует JSON в ShipAssembledEvent
func DecodeShipAssembledEvent(data []byte) (*model.ShipAssembledEvent, error) {
	var event model.ShipAssembledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
