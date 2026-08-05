package functions

import (
	"context"

	types2 "github.com/gorundebug/model/pkg/types"
	"github.com/gorundebug/orderservice/internal/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.MapFunction[*types2.OrderItemResult, *types.OrderState] = (*MapOrderItemResultToOrderState)(nil)

// MapOrderItemResultToOrderState
type MapOrderItemResultToOrderState struct{}

func (f *MapOrderItemResultToOrderState) Map(ctx context.Context, _ runtime.Stream, value *types2.OrderItemResult, out runtime.Collect[*types.OrderState]) {
	status := "PARTIALLY_CONFIRMED"
	if value.Reserved {
		status = "CONFIRMED"
	}
	out.Out(ctx, &types.OrderState{
		OrderID:        value.OrderID,
		Status:         status,
		ConfirmedItems: []*types2.OrderItemResult{value},
	})
}

// MakeMapOrderItemResultToOrderState is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeMapOrderItemResultToOrderState(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.MapStreamConfig) (*MapOrderItemResultToOrderState, error) {
	return &MapOrderItemResultToOrderState{}, nil
}
