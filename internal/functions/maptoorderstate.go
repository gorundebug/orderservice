package functions

import (
	"context"

	"github.com/gorundebug/orderservice/internal/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.MapFunction[*types.Order, *types.OrderState] = (*MapToOrderState)(nil)

// MapToOrderState
type MapToOrderState struct{}

func (f *MapToOrderState) Map(ctx context.Context, _ runtime.Stream, value *types.Order, out runtime.Collect[*types.OrderState]) {
	out.Out(ctx, &types.OrderState{
		OrderID:     value.ID,
		Status:      "TIMED_OUT",
		TotalAmount: value.TotalAmount,
	})
}

// MakeMapToOrderState is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeMapToOrderState(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.MapStreamConfig) (*MapToOrderState, error) {
	return &MapToOrderState{}, nil
}
