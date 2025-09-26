package inventory

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type InventoryServiceTestSuite struct {
	suite.Suite
}

func TestInventoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryServiceTestSuite))
}
