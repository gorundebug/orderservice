package order

import (
	"context"
	"testing"
	"time"

	"github.com/gorundebug/model/pkg/types"
	types2 "github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

func TestMapToOrderProcessed_Map(t *testing.T) {
	f := &MapToOrderProcessed{}
	var collected []*types.OrderProcessed
	out := runtime.CollectFunc[*types.OrderProcessed](func(_ context.Context, v *types.OrderProcessed) {
		collected = append(collected, v)
	})
	processedAt := time.Date(2026, time.August, 16, 12, 30, 0, 0, time.UTC)
	value := &types2.OrderState{
		OrderID: "order-123",
		Status:  "PARTIALLY_CONFIRMED",
		ConfirmedItems: []*types.OrderItemResult{
			{Reserved: true},
			nil,
			{Reserved: false},
		},
		ProcessedAt: processedAt,
	}
	f.Map(context.Background(), nil, value, out)
	if assert.Len(t, collected, 1) {
		assert.Equal(t, "order-123", collected[0].OrderID)
		assert.Equal(t, "PARTIALLY_CONFIRMED", collected[0].Status)
		assert.Equal(t, processedAt, collected[0].ProcessedAt)
		assert.Equal(t, 3, collected[0].TotalItems)
		assert.Equal(t, 1, collected[0].ConfirmedItems)
		assert.Equal(t, "PARTIALLY_CONFIRMED", collected[0].FailureReason)
	}
}
