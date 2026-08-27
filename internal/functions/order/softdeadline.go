package order

import (
	"context"

	"github.com/gorundebug/orderservice/internal/types"

	"time"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.DelayFunction[*types.Order] = (*SoftDeadline)(nil)

// SoftDeadline
type SoftDeadline struct{}

func (f *SoftDeadline) Duration(ctx context.Context, stream runtime.Stream, _ *types.Order) time.Duration {
	cfg := stream.GetConfig().(*runtimecfg.DelayStreamConfig)
	margin := time.Duration(cfg.Duration) * time.Millisecond
	deadline, ok := ctx.Deadline()
	if !ok {
		return margin
	}
	d := time.Until(deadline) - margin
	if d < 0 {
		return 0
	}
	return d
}

func (f *SoftDeadline) DelayError(_ context.Context, _ runtime.Stream, _ *types.Order, _ error, _ runtime.Collect[*types.Order]) {
	// drop: if the delay itself errors the order is already being handled by the timeout branch or is stale
}

// MakeSoftDeadline is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeSoftDeadline(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.DelayStreamConfig) (*SoftDeadline, error) {
	return &SoftDeadline{}, nil
}
