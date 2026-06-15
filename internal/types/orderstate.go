package types

// Processing result of an order. Fields: OrderID string, Status string (CONFIRMED — all items reserved;
// PARTIALLY_CONFIRMED — some items out of stock; TIMED_OUT — order timed out), ConfirmedItems []OrderItemResult,
// TotalAmount float64, ProcessedAt time.Time.
type OrderState struct {
}
