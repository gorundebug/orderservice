package order

import (
	"context"
	"testing"

	"github.com/gorundebug/model/pkg/types"
	types2 "github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Every item is emitted independently for inventory processing with the parent order ID assigned.

func TestProcessOrderItems_FlatMap(t *testing.T) {
	f := &ProcessOrderItems{}
	var collected []*types.OrderItem
	out := runtime.CollectFunc[*types.OrderItem](func(_ context.Context, v *types.OrderItem) {
		collected = append(collected, v)
	})
	value := &types2.Order{
		ID: "order-123",
		Items: []*types.OrderItem{
			{ItemID: "item-1", SKU: "SKU-001", Quantity: 2, UnitPrice: 10.0},
			{ItemID: "item-2", SKU: "SKU-002", Quantity: 1, UnitPrice: 5.0},
		},
	}
	f.FlatMap(context.Background(), nil, value, out)
	assert.Len(t, collected, 2)
	assert.Equal(t, "order-123", collected[0].OrderID)
	assert.Equal(t, "item-1", collected[0].ItemID)
	assert.Equal(t, "order-123", collected[1].OrderID)
	assert.Equal(t, "item-2", collected[1].ItemID)
}

func TestProcessOrderItems_FlatMap_EmptyItems(t *testing.T) {
	f := &ProcessOrderItems{}
	var collected []*types.OrderItem
	out := runtime.CollectFunc[*types.OrderItem](func(_ context.Context, v *types.OrderItem) {
		collected = append(collected, v)
	})
	value := &types2.Order{ID: "order-empty", Items: nil}
	f.FlatMap(context.Background(), nil, value, out)
	assert.Empty(t, collected)
}
