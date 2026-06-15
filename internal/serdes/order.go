package serdes

import (
	"fmt"

	"github.com/gorundebug/orderservice/internal/types"
)

type OrderSerde struct{}

func (s *OrderSerde) IsStub() bool {
	return false
}

func (s *OrderSerde) SerializeObj(value interface{}, b []byte) ([]byte, error) {
	v, ok := value.(*types.Order)
	if !ok {
		return nil, fmt.Errorf("value is not *types.Order")
	}
	return s.Serialize(v, b)
}

func (s *OrderSerde) DeserializeObj(data []byte) (interface{}, error) {
	return s.Deserialize(data)
}

func (s *OrderSerde) Serialize(value *types.Order, b []byte) ([]byte, error) {
	// TODO: Need to be implemented
	return nil, fmt.Errorf("serialize method for the 'OrderSerde' class is not implemented")
}

func (s *OrderSerde) Deserialize(data []byte) (*types.Order, error) {
	// TODO: Need to be implemented
	return nil, fmt.Errorf("deserialize method for the 'OrderSerde' class is not implemented")
}
