package types

import (
	"time"

	modeltypes "github.com/gorundebug/model_go/pkg/types"
)

// E-commerce order submitted by a customer. Fields: ID string (UUID, set at HTTP ingress), CustomerID string, Items
// []OrderItem, TotalAmount float64, CreatedAt time.Time, TraceID string (propagated from X-Request-ID header).
type Order struct {
	ID          string
	CustomerID  string
	Items       []*modeltypes.OrderItem
	TotalAmount float64
	CreatedAt   time.Time
	TraceID     string
}
