package functions

import (
	"context"
	"testing"

	"github.com/gorundebug/model/pkg/types"
	types2 "github.com/gorundebug/orderservice/internal/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Expand an Order into individual OrderItem messages — one sc.Collect call per element of Order.Items. Copy Order.ID
// into each emitted OrderItem.OrderID.

func TestProcessOrderItems_FlatMap(t *testing.T) {
	t.Skip("not yet implemented") // TODO: remove when implementation is ready
	f := &ProcessOrderItems{}
	var collected []*types.OrderItem
	out := runtime.CollectFunc[*types.OrderItem](func(_ context.Context, v *types.OrderItem) {
		collected = append(collected, v)
	})
	var value *types2.Order
	// TODO: populate value with meaningful test data
	f.FlatMap(context.Background(), nil, value, out)
	assert.NotEmpty(t, collected)
}
