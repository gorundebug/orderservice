package types

import (
	"time"

	modeltypes "github.com/gorundebug/model/pkg/types"
)

// Current processing state of an order, returned to the HTTP caller. Fields: OrderID string, Status string (CONFIRMED
// — all items reserved; PARTIALLY_CONFIRMED — some items out of stock; TIMED_OUT — soft deadline expired before
// all items were processed), ConfirmedItems []OrderItemResult, TotalAmount float64, ProcessedAt time.Time.
type OrderState struct {
	OrderID        string
	Status         string
	ConfirmedItems []*modeltypes.OrderItemResult
	TotalAmount    float64
	ProcessedAt    time.Time
}
