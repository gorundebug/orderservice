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

var _ transformation.MapFunction[*types2.OrderState, *types.OrderProcessed] = (*MapToOrderProcessed)(nil)

// MapToOrderProcessed
type MapToOrderProcessed struct{}

func (f *MapToOrderProcessed) Map(ctx context.Context, _ runtime.Stream, value *types2.OrderState, out runtime.Collect[*types.OrderProcessed]) {
	confirmedItems := 0
	for _, item := range value.ConfirmedItems {
		if item != nil && item.Reserved {
			confirmedItems++
		}
	}
	failureReason := ""
	if value.Status != "CONFIRMED" {
		failureReason = value.Status
	}
	out.Out(ctx, &types.OrderProcessed{
		OrderID:        value.OrderID,
		Status:         value.Status,
		ProcessedAt:    value.ProcessedAt,
		TotalItems:     len(value.ConfirmedItems),
		ConfirmedItems: confirmedItems,
		FailureReason:  failureReason,
	})
}

// MakeMapToOrderProcessed is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeMapToOrderProcessed(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.MapStreamConfig) (*MapToOrderProcessed, error) {
	return &MapToOrderProcessed{}, nil
}
