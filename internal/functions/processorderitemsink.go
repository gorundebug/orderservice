package functions

import (
	"context"

	datasinkgrpc "github.com/gorundebug/servicelib/datasink/grpc"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"

	"github.com/gorundebug/inventory_service_api/pkg/generated/proto/inventoryserviceapi/processorderitem"
	"github.com/gorundebug/model/pkg/types"
	types2 "github.com/gorundebug/orderservice/internal/types"
)

// processOrderItemSinkHandler is a type alias for the EndpointHandler generic instantiation used throughout this file.
type processOrderItemSinkHandler = datasinkgrpc.EndpointHandler[ProcessOrderItemSinkHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, *types2.OrderState]

var _ processOrderItemSinkHandler = (*ProcessOrderItemSink)(nil)

func MakeEndpointConsumerProcessOrderItemSink(stream runtime.TypedSinkStreamWithResult[*types.OrderItem, *types.OrderItemResult, *types2.OrderState], handler processOrderItemSinkHandler, callFn datasinkgrpc.NoStreamingClientFunction[*processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse]) (runtime.Consumer[*types.OrderItem], error) {
	return datasinkgrpc.MakeGRPCNoStreamingEndpointConsumer[ProcessOrderItemSinkHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, *types2.OrderState](stream, handler, callFn)
}

type processOrderItemSinkData struct {
	orderID      string
	itemID       string
	sku          string
	requestedQty int
	unitPrice    float64
}

// ProcessOrderItemSinkHandlerState holds per-request state created by BeginRequest for each outgoing gRPC call.
type ProcessOrderItemSinkHandlerState struct {
	data *processOrderItemSinkData
}

// ProcessOrderItemSink
type ProcessOrderItemSink struct{}

func (ep *ProcessOrderItemSink) BeginRequest(ctx context.Context, _ datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, *types2.OrderState]) (context.Context, ProcessOrderItemSinkHandlerState, error) {
	return ctx, ProcessOrderItemSinkHandlerState{data: &processOrderItemSinkData{}}, nil
}

func (ep *ProcessOrderItemSink) ConsumeMessage(ctx context.Context, _ datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, *types2.OrderState], handlerState ProcessOrderItemSinkHandlerState, value *types.OrderItem, sender datasinkgrpc.Sender[*processorderitem.ProcessOrderItemRequest], resultCtx datasinkgrpc.ResultContext) error {
	handlerState.data.orderID = value.OrderID
	handlerState.data.itemID = value.ItemID
	handlerState.data.sku = value.SKU
	handlerState.data.requestedQty = value.Quantity
	handlerState.data.unitPrice = value.UnitPrice

	req := &processorderitem.ProcessOrderItemRequest{
		OrderId:  value.OrderID,
		ItemId:   value.ItemID,
		Sku:      value.SKU,
		Quantity: int32(value.Quantity),
	}
	if err := sender.Send(ctx, req); err != nil {
		return err
	}
	return nil
}

func (ep *ProcessOrderItemSink) HandleResponse(ctx context.Context, sc datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, *types2.OrderState], handlerState ProcessOrderItemSinkHandlerState, response *processorderitem.ProcessOrderItemResponse) error {
	result := &types.OrderItemResult{
		OrderID:      handlerState.data.orderID,
		ItemID:       handlerState.data.itemID,
		SKU:          handlerState.data.sku,
		RequestedQty: handlerState.data.requestedQty,
		AvailableQty: int(response.AvailableQty),
		Reserved:     response.Reserved,
		Status:       response.Status,
		UnitPrice:    handlerState.data.unitPrice,
	}
	sc.Collect(ctx, result)
	return nil
}

func (ep *ProcessOrderItemSink) EndRequest(ctx context.Context, sc datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, *types2.OrderState], err error, handlerState ProcessOrderItemSinkHandlerState) {
	if err == nil {
		return
	}
	sc.Collect(ctx, &types.OrderItemResult{
		OrderID:      handlerState.data.orderID,
		ItemID:       handlerState.data.itemID,
		SKU:          handlerState.data.sku,
		RequestedQty: handlerState.data.requestedQty,
		Reserved:     false,
		Status:       "PROCESSING_ERROR",
		UnitPrice:    handlerState.data.unitPrice,
		Error:        err.Error(),
	})
}

// MakeProcessOrderItemSink implements the handler for the ProcessOrderItem gRPC sink endpoint.
func MakeProcessOrderItemSink(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.GrpcEndpointConfig) (*ProcessOrderItemSink, error) {
	return &ProcessOrderItemSink{}, nil
}
