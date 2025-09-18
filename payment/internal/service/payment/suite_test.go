package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
}

func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}
