package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/payment/internal/model"
	repomodel "github.com/nimbodex/microservices-factory/payment/internal/repository/model"
)

// ToRepoPayment converts service model to repository model
func ToRepoPayment(servicePayment *model.Payment) *repomodel.Payment {
	if servicePayment == nil {
		return nil
	}

	return &repomodel.Payment{
		UUID:            servicePayment.UUID.String(),
		OrderUUID:       servicePayment.OrderUUID.String(),
		PaymentMethod:   string(servicePayment.PaymentMethod),
		Amount:          servicePayment.Amount,
		Status:          string(servicePayment.Status),
		TransactionUUID: servicePayment.TransactionUUID.String(),
		CreatedAt:       servicePayment.CreatedAt,
		UpdatedAt:       servicePayment.UpdatedAt,
	}
}

// FromRepoPayment converts repository model to service model
func FromRepoPayment(repoPayment *repomodel.Payment) (*model.Payment, error) {
	if repoPayment == nil {
		return nil, fmt.Errorf("repoPayment cannot be nil")
	}

	paymentUUID, err := uuid.Parse(repoPayment.UUID)
	if err != nil {
		return nil, err
	}

	orderUUID, err := uuid.Parse(repoPayment.OrderUUID)
	if err != nil {
		return nil, err
	}

	transactionUUID, err := uuid.Parse(repoPayment.TransactionUUID)
	if err != nil {
		return nil, err
	}

	return &model.Payment{
		UUID:            paymentUUID,
		OrderUUID:       orderUUID,
		PaymentMethod:   model.PaymentMethod(repoPayment.PaymentMethod),
		Amount:          repoPayment.Amount,
		Status:          model.PaymentStatus(repoPayment.Status),
		TransactionUUID: transactionUUID,
		CreatedAt:       repoPayment.CreatedAt,
		UpdatedAt:       repoPayment.UpdatedAt,
	}, nil
}
