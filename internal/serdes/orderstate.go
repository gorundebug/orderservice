package serdes

import (
	"fmt"

	"github.com/gorundebug/orderservice/internal/types"
)

type OrderStateSerde struct{}

func (s *OrderStateSerde) IsStub() bool {
	return false
}

func (s *OrderStateSerde) SerializeObj(value interface{}, b []byte) ([]byte, error) {
	v, ok := value.(*types.OrderState)
	if !ok {
		return nil, fmt.Errorf("value is not *types.OrderState")
	}
	return s.Serialize(v, b)
}

func (s *OrderStateSerde) DeserializeObj(data []byte) (interface{}, error) {
	return s.Deserialize(data)
}

func (s *OrderStateSerde) Serialize(value *types.OrderState, b []byte) ([]byte, error) {
	// TODO: Need to be implemented
	return nil, fmt.Errorf("serialize method for the 'OrderStateSerde' class is not implemented")
}

func (s *OrderStateSerde) Deserialize(data []byte) (*types.OrderState, error) {
	// TODO: Need to be implemented
	return nil, fmt.Errorf("deserialize method for the 'OrderStateSerde' class is not implemented")
}
