package serdes

import (
	"fmt"

	"github.com/gorundebug/inventoryservice/internal/types"
)

type InventoryFailureSerde struct{}

func (s *InventoryFailureSerde) IsStub() bool {
	return false
}

func (s *InventoryFailureSerde) SerializeObj(value interface{}, b []byte) ([]byte, error) {
	v, ok := value.(*types.InventoryFailure)
	if !ok {
		return nil, fmt.Errorf("value is not *types.InventoryFailure")
	}
	return s.Serialize(v, b)
}

func (s *InventoryFailureSerde) DeserializeObj(data []byte) (interface{}, error) {
	return s.Deserialize(data)
}

func (s *InventoryFailureSerde) Serialize(value *types.InventoryFailure, b []byte) ([]byte, error) {
	// TODO: Need to be implemented
	return nil, fmt.Errorf("serialize method for the 'InventoryFailureSerde' class is not implemented")
}

func (s *InventoryFailureSerde) Deserialize(data []byte) (*types.InventoryFailure, error) {
	// TODO: Need to be implemented
	return nil, fmt.Errorf("deserialize method for the 'InventoryFailureSerde' class is not implemented")
}
