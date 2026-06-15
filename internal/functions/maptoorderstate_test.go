package functions

import (
	"context"
	"testing"

	"github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Convert an Order that reached the soft deadline into an OrderState. Set OrderID from Order.ID; set Status to
// TIMED_OUT; leave ConfirmedItems nil.

func TestMapToOrderState_Map(t *testing.T) {
	t.Skip("not yet implemented") // TODO: remove when implementation is ready
	f := &MapToOrderState{}
	var collected []*types.OrderState
	out := runtime.CollectFunc[*types.OrderState](func(_ context.Context, v *types.OrderState) {
		collected = append(collected, v)
	})
	var value *types.Order
	// TODO: populate value with meaningful test data
	f.Map(context.Background(), nil, value, out)
	assert.NotEmpty(t, collected)
}
