package functions

import (
	"context"
	"testing"

	"github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Convert an Order that reached the soft deadline (timeout branch) into a partial OrderState{Status: TIMED_OUT}. Copy
// OrderID and whichever fields are available; leave ConfirmedItems empty or with the items collected so far. This
// response unblocks the waiting HTTP client when inventory checks have not finished within the allowed window.

func TestMapToOrderState_Map(t *testing.T) {
	f := &MapToOrderState{}
	var collected []*types.OrderState
	out := runtime.CollectFunc[*types.OrderState](func(_ context.Context, v *types.OrderState) {
		collected = append(collected, v)
	})
	value := &types.Order{
		ID:          "order-789",
		CustomerID:  "cust-1",
		TotalAmount: 150.0,
	}
	f.Map(context.Background(), nil, value, out)
	assert.Len(t, collected, 1)
	assert.Equal(t, "order-789", collected[0].OrderID)
	assert.Equal(t, "TIMED_OUT", collected[0].Status)
	assert.Equal(t, 150.0, collected[0].TotalAmount)
	assert.Empty(t, collected[0].ConfirmedItems)
}
