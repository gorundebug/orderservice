package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	modeltypes "github.com/gorundebug/model/pkg/types"
	datasourcehttp "github.com/gorundebug/servicelib/datasource/http"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"

	"github.com/gorundebug/order_service_api/pkg/generated/openapi/orderserviceapi/processorder"
	"github.com/gorundebug/orderservice/internal/types"
)

// ProcessOrderSourceType is the typed handler function for the ProcessOrderSource HTTP endpoint.
type ProcessOrderSourceType = datasourcehttp.HTTPHandler

// processOrderHandler is a type alias for the EndpointHandler generic instantiation used throughout this file.
type processOrderHandler = datasourcehttp.EndpointHandler[ProcessOrderSourceHandlerState, *processorder.ProcessOrderRequest, *processorder.ProcessOrderResponse, *types.Order, *types.OrderState, error]

var _ processOrderHandler = (*ProcessOrderSource)(nil)

func MakeEndpointConsumerProcessOrderSource(stream runtime.TypedInputStream[*types.Order, *types.OrderState, error], handler processOrderHandler) (runtime.Consumer[*types.Order], ProcessOrderSourceType, error) {
	return datasourcehttp.MakeNetHTTPEndpointConsumer[ProcessOrderSourceHandlerState, *processorder.ProcessOrderRequest, *processorder.ProcessOrderResponse, *types.Order, *types.OrderState, error](stream, handler)
}

// ProcessOrderSourceHandlerState holds per-request state created by BeginRequest for each incoming HTTP request.
type ProcessOrderSourceHandlerState struct {
	req    *processorder.ProcessOrderRequest
	cancel context.CancelFunc
}

// ProcessOrderSource
type ProcessOrderSource struct {
	timeout time.Duration
}

func (ep *ProcessOrderSource) BeginRequest(ctx context.Context, _ datasourcehttp.StreamContext[*types.Order, *types.OrderState, error], data datasourcehttp.HandlerData) (context.Context, ProcessOrderSourceHandlerState, error) {
	var req processorder.ProcessOrderRequest
	if err := json.NewDecoder(data.Request.Body).Decode(&req); err != nil {
		http.Error(data.Writer, "invalid JSON body", http.StatusBadRequest)
		return ctx, ProcessOrderSourceHandlerState{}, err
	}
	if len(req.Items) == 0 {
		err := errors.New("items must not be empty")
		http.Error(data.Writer, err.Error(), http.StatusBadRequest)
		return ctx, ProcessOrderSourceHandlerState{}, err
	}
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			err := errors.New("all quantities must be positive")
			http.Error(data.Writer, err.Error(), http.StatusBadRequest)
			return ctx, ProcessOrderSourceHandlerState{}, err
		}
	}
	ctx, cancel := context.WithTimeout(ctx, ep.timeout)
	return ctx, ProcessOrderSourceHandlerState{req: &req, cancel: cancel}, nil
}

func (ep *ProcessOrderSource) ConsumeMessage(ctx context.Context, sc datasourcehttp.StreamContext[*types.Order, *types.OrderState, error], handlerState ProcessOrderSourceHandlerState, data datasourcehttp.HandlerData, resultCtx datasourcehttp.ResultContext[ProcessOrderSourceHandlerState, *processorder.ProcessOrderRequest, *processorder.ProcessOrderResponse, *types.Order, *types.OrderState, error]) error {
	orderID := data.Request.Header.Get("X-Request-ID")
	if orderID == "" {
		orderID = uuid.New().String()
	}

	req := handlerState.req
	items := make([]*modeltypes.OrderItem, 0, len(req.Items))
	var total float64
	for _, ri := range req.Items {
		item := &modeltypes.OrderItem{
			OrderID:  orderID,
			ItemID:   ri.ItemId,
			SKU:      ri.Sku,
			Quantity: int(ri.Quantity),
		}
		if ri.UnitPrice != nil {
			item.UnitPrice = *ri.UnitPrice
			total += float64(ri.Quantity) * *ri.UnitPrice
		}
		items = append(items, item)
	}

	customerID := ""
	if req.CustomerId != nil {
		customerID = *req.CustomerId
	}

	order := &types.Order{
		ID:          orderID,
		CustomerID:  customerID,
		Items:       items,
		TotalAmount: total,
		CreatedAt:   time.Now(),
		TraceID:     data.Request.Header.Get("X-Trace"),
	}

	var mu sync.Mutex
	results := make([]*modeltypes.OrderItemResult, 0, len(items))
	responseSent := false
	resultCtx.SetResultCallback(orderID, func(ctx context.Context, sc datasourcehttp.StreamContext[*types.Order, *types.OrderState, error], _ ProcessOrderSourceHandlerState, state *types.OrderState, data datasourcehttp.HandlerData) bool {
		mu.Lock()
		defer mu.Unlock()
		if responseSent {
			return true
		}

		if state.Status != "TIMED_OUT" {
			results = append(results, state.ConfirmedItems...)
			if len(results) < len(items) {
				return false
			}
		}

		status := state.Status
		if status != "TIMED_OUT" {
			status = "CONFIRMED"
			for _, result := range results {
				if !result.Reserved {
					status = "PARTIALLY_CONFIRMED"
					break
				}
			}
		}

		totalAmount := 0.0
		for _, result := range results {
			totalAmount += result.UnitPrice * float64(result.RequestedQty)
		}
		if len(results) == 0 {
			totalAmount = order.TotalAmount
		}

		resp := buildProcessOrderSourceResponse(&types.OrderState{
			OrderID:        order.ID,
			Status:         status,
			ConfirmedItems: results,
			TotalAmount:    totalAmount,
			ProcessedAt:    time.Now(),
		})
		data.Writer.Header().Set("Content-Type", "application/json")
		data.Writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(data.Writer).Encode(resp); err != nil {
			return false
		}
		responseSent = true
		resultCtx.Done()
		return true
	})

	sc.Collect(ctx, order)
	return nil
}

func (ep *ProcessOrderSource) GetMessageID(_ context.Context, _ datasourcehttp.StreamContext[*types.Order, *types.OrderState, error], _ ProcessOrderSourceHandlerState, value *types.OrderState) string {
	return value.OrderID
}

func (ep *ProcessOrderSource) EndRequest(_ context.Context, _ datasourcehttp.StreamContext[*types.Order, *types.OrderState, error], err error, handlerState ProcessOrderSourceHandlerState, data datasourcehttp.HandlerData) {
	handlerState.cancel()
	if err != nil {
		http.Error(data.Writer, "internal server error", http.StatusInternalServerError)
	}
}

func buildProcessOrderSourceResponse(state *types.OrderState) *processorder.ProcessOrderResponse {
	resp := &processorder.ProcessOrderResponse{
		OrderId:     &state.OrderID,
		Status:      &state.Status,
		TotalAmount: &state.TotalAmount,
		ProcessedAt: &state.ProcessedAt,
	}
	if len(state.ConfirmedItems) > 0 {
		items := make([]processorder.ProcessOrderResponseItem, 0, len(state.ConfirmedItems))
		for _, ci := range state.ConfirmedItems {
			avail := int32(ci.AvailableQty)
			item := processorder.ProcessOrderResponseItem{
				ItemId:       &ci.ItemID,
				Sku:          &ci.SKU,
				AvailableQty: &avail,
				Reserved:     &ci.Reserved,
				Status:       &ci.Status,
			}
			if ci.Error != "" {
				item.Error = &ci.Error
			}
			items = append(items, item)
		}
		resp.ConfirmedItems = &items
	}
	return resp
}

// MakeProcessOrderSource implements the handler for the ProcessOrderSource HTTP source endpoint.
func MakeProcessOrderSource(_ context.Context, _ environment.ServiceEnvironment, cfg *runtimecfg.HttpEndpointConfig) (*ProcessOrderSource, error) {
	const defaultTimeout = 5 * time.Second
	timeout := defaultTimeout
	if v := cfg.GetProperty("timeout"); v != nil {
		switch ms := v.(type) {
		case int:
			timeout = time.Duration(ms) * time.Millisecond
		case float64:
			timeout = time.Duration(ms) * time.Millisecond
		}
	}
	return &ProcessOrderSource{timeout: timeout}, nil
}
