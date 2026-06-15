package functions

import (
	"context"
	"testing"
	"time"

	"github.com/gorundebug/orderservice/internal/types"
	"github.com/stretchr/testify/assert"
)

// Cast stream.GetConfig() to *runtimecfg.DelayStreamConfig and convert cfg.Duration (int, milliseconds) to
// time.Duration — this is the safety margin. If ctx has no deadline (ctx.Deadline() ok==false), return the margin
// directly. Otherwise compute time.Until(deadline) minus the margin: if the result is negative return 0, otherwise
// return it.

func TestSoftDeadline_Duration(t *testing.T) {
	t.Skip("not yet implemented") // TODO: remove when implementation is ready
	f := &SoftDeadline{}
	var value *types.Order
	// TODO: populate value with meaningful test data.
	// Note: if your implementation accesses stream.GetConfig(), pass a mock runtime.Stream instead of nil.
	result := f.Duration(context.Background(), nil, value)
	assert.GreaterOrEqual(t, result, time.Duration(0))
}

func TestSoftDeadline_DelayError(t *testing.T) {
	t.Skip("not yet implemented") // TODO: remove when implementation is ready
	f := &SoftDeadline{}
	var value *types.Order
	// TODO: populate value with meaningful test data and verify whether the value is re-emitted or dropped.
	f.DelayError(context.Background(), nil, value, context.DeadlineExceeded, nil)
}
