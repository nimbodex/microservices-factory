package converter

import (
	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/order/internal/model"
	repomodel "github.com/nimbodex/microservices-factory/order/internal/repository/model"
)

// ToRepoOrder converts service model to repository model
func ToRepoOrder(order *model.Order) *repomodel.Order {
	partUUIDs := make([]string, len(order.PartUUIDs))
	for i, uuid := range order.PartUUIDs {
		partUUIDs[i] = uuid.String()
	}

	return &repomodel.Order{
		UUID:      order.UUID.String(),
		UserUUID:  order.UserUUID.String(),
		PartUUIDs: partUUIDs,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

// FromRepoOrder converts repository model to service model
func FromRepoOrder(repoOrder *repomodel.Order) (*model.Order, error) {
	orderUUID, err := uuid.Parse(repoOrder.UUID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(repoOrder.UserUUID)
	if err != nil {
		return nil, err
	}

	partUUIDs := make([]uuid.UUID, len(repoOrder.PartUUIDs))
	for i, uuidStr := range repoOrder.PartUUIDs {
		partUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, err
		}
		partUUIDs[i] = partUUID
	}

	return &model.Order{
		UUID:      orderUUID,
		UserUUID:  userUUID,
		PartUUIDs: partUUIDs,
		Status:    model.OrderStatus(repoOrder.Status),
		CreatedAt: repoOrder.CreatedAt,
		UpdatedAt: repoOrder.UpdatedAt,
	}, nil
}

// ToRepoPart converts service model to repository model
func ToRepoPart(part *model.Part) *repomodel.Part {
	return &repomodel.Part{
		UUID:  part.UUID.String(),
		Name:  part.Name,
		Price: part.Price,
	}
}

// FromRepoPart converts repository model to service model
func FromRepoPart(repoPart *repomodel.Part) (*model.Part, error) {
	partUUID, err := uuid.Parse(repoPart.UUID)
	if err != nil {
		return nil, err
	}

	return &model.Part{
		UUID:  partUUID,
		Name:  repoPart.Name,
		Price: repoPart.Price,
	}, nil
}
