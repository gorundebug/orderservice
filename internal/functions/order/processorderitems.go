package order

import (
	"context"

	"github.com/gorundebug/model/pkg/types"
	types2 "github.com/gorundebug/orderservice/internal/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.FlatMapFunction[*types2.Order, *types.OrderItem] = (*ProcessOrderItems)(nil)

// ProcessOrderItems
type ProcessOrderItems struct{}

func (f *ProcessOrderItems) FlatMap(ctx context.Context, _ runtime.Stream, value *types2.Order, out runtime.Collect[*types.OrderItem]) {
	for _, item := range value.Items {
		itemCopy := *item
		itemCopy.OrderID = value.ID
		out.Out(ctx, &itemCopy)
	}
}

// MakeProcessOrderItems is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeProcessOrderItems(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.FlatMapStreamConfig) (*ProcessOrderItems, error) {
	return &ProcessOrderItems{}, nil
}
