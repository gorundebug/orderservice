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

func (f *MapToOrderState) Map(_ context.Context, _ runtime.Stream, value *types.Order, out runtime.Collect[*types.OrderState]) {
	//TODO: Need to be implemented
	// Convert an Order that reached the soft deadline into an OrderState. Set OrderID from Order.ID; set Status to
	// TIMED_OUT; leave ConfirmedItems nil.
}

// MakeMapToOrderState is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeMapToOrderState(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.MapStreamConfig) (*MapToOrderState, error) {
	return &MapToOrderState{}, nil
}
