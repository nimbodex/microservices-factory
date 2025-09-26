package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	suite.Suite
}

func TestOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
