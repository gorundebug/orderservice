package serdes

import (
	"encoding/json"
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

func (s *OrderStateSerde) Serialize(value *types.OrderState, _ []byte) ([]byte, error) {
	return json.Marshal(value)
}

func (s *OrderStateSerde) Deserialize(data []byte) (*types.OrderState, error) {
	var v types.OrderState
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
