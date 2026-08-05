package functions

import (
	"context"
	"testing"
	"time"

	"github.com/gorundebug/orderservice/internal/types"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/stretchr/testify/assert"
)

// Returns the soft-deadline duration to cut off from the request context deadline. cfg.Duration (milliseconds) is the
// safety margin reserved for assembling and returning a partial result. When the soft deadline fires, the Order flows
// through the timeout branch so the service can collect whichever OrderItems have already been processed and return at
// least a partial OrderState to the caller, rather than waiting indefinitely for all items to complete.

type mockDelayStream struct {
	cfg *runtimecfg.DelayStreamConfig
}

func (m *mockDelayStream) GetName() string                                { return "test" }
func (m *mockDelayStream) GetTransformationName() string                  { return "Delay" }
func (m *mockDelayStream) GetTypeName() string                            { return "Order" }
func (m *mockDelayStream) GetID() int                                     { return 0 }
func (m *mockDelayStream) GetConfig() runtimecfg.StreamConfig             { return m.cfg }
func (m *mockDelayStream) GetEnvironment() environment.ServiceEnvironment { return nil }

func TestSoftDeadline_Duration_WithDeadline(t *testing.T) {
	f := &SoftDeadline{}
	value := &types.Order{ID: "order-1"}
	stream := &mockDelayStream{cfg: &runtimecfg.DelayStreamConfig{Duration: 200}}

	// context deadline 1 second from now, margin 200ms → expected ~800ms
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result := f.Duration(ctx, stream, value)
	assert.GreaterOrEqual(t, result, 700*time.Millisecond)
	assert.LessOrEqual(t, result, 900*time.Millisecond)
}

func TestSoftDeadline_Duration_NoDeadline(t *testing.T) {
	f := &SoftDeadline{}
	value := &types.Order{ID: "order-2"}
	stream := &mockDelayStream{cfg: &runtimecfg.DelayStreamConfig{Duration: 300}}

	// no deadline → fall back to margin
	result := f.Duration(context.Background(), stream, value)
	assert.Equal(t, 300*time.Millisecond, result)
}

func TestSoftDeadline_Duration_MarginExceedsDeadline(t *testing.T) {
	f := &SoftDeadline{}
	value := &types.Order{ID: "order-3"}
	stream := &mockDelayStream{cfg: &runtimecfg.DelayStreamConfig{Duration: 5000}}

	// deadline only 100ms away but margin is 5s → clamp to 0
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result := f.Duration(ctx, stream, value)
	assert.Equal(t, time.Duration(0), result)
}

func TestSoftDeadline_DelayError(t *testing.T) {
	f := &SoftDeadline{}
	value := &types.Order{ID: "order-4"}
	// DelayError should drop the value without panicking
	f.DelayError(context.Background(), nil, value, context.DeadlineExceeded, nil)
}
