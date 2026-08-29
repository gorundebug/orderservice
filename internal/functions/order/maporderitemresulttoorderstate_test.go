package order

import (
	"context"
	"testing"

	types2 "github.com/gorundebug/model_go/pkg/types"
	"github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// A single inventory result becomes an order result with the same order ID. Reserved items are CONFIRMED;
// all other outcomes are PARTIALLY_CONFIRMED.

func TestMapOrderItemResultToOrderState_Map_Confirmed(t *testing.T) {
	f := &MapOrderItemResultToOrderState{}
	var collected []*types.OrderState
	out := runtime.CollectFunc[*types.OrderState](func(_ context.Context, v *types.OrderState) {
		collected = append(collected, v)
	})
	value := &types2.OrderItemResult{
		OrderID:      "order-123",
		ItemID:       "item-1",
		SKU:          "SKU-001",
		RequestedQty: 2,
		AvailableQty: 2,
		Reserved:     true,
		Status:       "CONFIRMED",
	}
	f.Map(context.Background(), nil, value, out)
	assert.Len(t, collected, 1)
	assert.Equal(t, "order-123", collected[0].OrderID)
	assert.Equal(t, "CONFIRMED", collected[0].Status)
	assert.Len(t, collected[0].ConfirmedItems, 1)
}

func TestMapOrderItemResultToOrderState_Map_PartiallyConfirmed(t *testing.T) {
	f := &MapOrderItemResultToOrderState{}
	var collected []*types.OrderState
	out := runtime.CollectFunc[*types.OrderState](func(_ context.Context, v *types.OrderState) {
		collected = append(collected, v)
	})
	value := &types2.OrderItemResult{
		OrderID:      "order-456",
		ItemID:       "item-2",
		SKU:          "SKU-002",
		RequestedQty: 5,
		AvailableQty: 3,
		Reserved:     false,
		Status:       "OUT_OF_STOCK",
	}
	f.Map(context.Background(), nil, value, out)
	assert.Len(t, collected, 1)
	assert.Equal(t, "order-456", collected[0].OrderID)
	assert.Equal(t, "PARTIALLY_CONFIRMED", collected[0].Status)
}
