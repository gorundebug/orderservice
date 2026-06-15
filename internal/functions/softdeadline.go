package functions

import (
	"context"

	"github.com/gorundebug/orderservice/internal/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
	"time"
)

var _ transformation.DelayFunction[*types.Order] = (*SoftDeadline)(nil)

// SoftDeadline
type SoftDeadline struct{}

func (f *SoftDeadline) Duration(_ context.Context, stream runtime.Stream, value *types.Order) time.Duration {
	//TODO: Need to be implemented
	// Cast stream.GetConfig() to *runtimecfg.DelayStreamConfig and convert cfg.Duration (int, milliseconds) to
	// time.Duration — this is the safety margin. If ctx has no deadline (ctx.Deadline() ok==false), return the margin
	// directly. Otherwise compute time.Until(deadline) minus the margin: if the result is negative return 0, otherwise
	// return it.

	// Runtime invariant: Delay streams always have *runtimecfg.DelayStreamConfig.
	// Therefore the type assertion cannot fail.
	cfg := stream.GetConfig().(*runtimecfg.DelayStreamConfig)
	return time.Duration(cfg.Duration) * time.Millisecond
}

func (f *SoftDeadline) DelayError(_ context.Context, _ runtime.Stream, _ *types.Order, _ error, _ runtime.Collect[*types.Order]) {
	//TODO: Need to be implemented — decide whether to re-emit the value or drop it
}

// MakeSoftDeadline is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeSoftDeadline(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.DelayStreamConfig) (*SoftDeadline, error) {
	return &SoftDeadline{}, nil
}
