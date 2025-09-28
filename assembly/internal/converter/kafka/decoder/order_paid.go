package decoder

import (
	"encoding/json"

	"github.com/nimbodex/microservices-factory/assembly/internal/model"
)

// DecodeOrderPaidEvent декодирует JSON в OrderPaidEvent
func DecodeOrderPaidEvent(data []byte) (*model.OrderPaidEvent, error) {
	var event model.OrderPaidEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
