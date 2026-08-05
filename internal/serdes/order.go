package serdes

import (
	"encoding/json"
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

func (s *OrderSerde) Serialize(value *types.Order, _ []byte) ([]byte, error) {
	return json.Marshal(value)
}

func (s *OrderSerde) Deserialize(data []byte) (*types.Order, error) {
	var v types.Order
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
