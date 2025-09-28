package kafka

import (
	"encoding/json"

	"github.com/nimbodex/microservices-factory/order/internal/model"
)

// EncodeOrderPaidEvent кодирует OrderPaidEvent в JSON
func EncodeOrderPaidEvent(event *model.OrderPaidEvent) ([]byte, error) {
	return json.Marshal(event)
}
