package functions

import (
	"context"
	"testing"

	types2 "github.com/gorundebug/model/pkg/types"
	"github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Convert a single OrderItemResult into an OrderState. Set OrderID from result.OrderID; set Status=CONFIRMED if
// result.Reserved==true, otherwise PARTIALLY_CONFIRMED. Set ConfirmedItems to a single-element slice containing result.

func TestMapOrderItemResultToOrderState_Map(t *testing.T) {
	t.Skip("not yet implemented") // TODO: remove when implementation is ready
	f := &MapOrderItemResultToOrderState{}
	var collected []*types.OrderState
	out := runtime.CollectFunc[*types.OrderState](func(_ context.Context, v *types.OrderState) {
		collected = append(collected, v)
	})
	var value *types2.OrderItemResult
	// TODO: populate value with meaningful test data
	f.Map(context.Background(), nil, value, out)
	assert.NotEmpty(t, collected)
}
