package functions

import (
	"context"

	datasinkgrpc "github.com/gorundebug/servicelib/datasink/grpc"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"

	"github.com/gorundebug/inventory_service_api/pkg/generated/proto/inventoryserviceapi/processorderitem"
	"github.com/gorundebug/model/pkg/types"
)

// processOrderItemSinkHandler is a type alias for the EndpointHandler generic instantiation used throughout this file.
type processOrderItemSinkHandler = datasinkgrpc.EndpointHandler[ProcessOrderItemSinkHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, error]

var _ processOrderItemSinkHandler = (*ProcessOrderItemSink)(nil)

func MakeEndpointConsumerProcessOrderItemSink(stream runtime.TypedSinkStreamWithResult[*types.OrderItem, *types.OrderItemResult, error], handler processOrderItemSinkHandler, callFn datasinkgrpc.NoStreamingClientFunction[*processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse]) (runtime.Consumer[*types.OrderItem], error) {
	return datasinkgrpc.MakeGRPCNoStreamingEndpointConsumer[ProcessOrderItemSinkHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, error](stream, handler, callFn)
}

// ProcessOrderItemSinkHandlerState holds per-request state created by BeginRequest for each outgoing gRPC call.
// Enables safe concurrent processing — no synchronization needed between requests.
// Add fields here to carry data across BeginRequest → ConsumeMessage → HandleResponse → EndRequest.
type ProcessOrderItemSinkHandlerState struct {
}

// ProcessOrderItemSink
// Outgoing unary gRPC call to the Inventory Service. [ConsumeMessage] map OrderItem → ProcessOrderItemRequest
// (OrderID, ItemID, SKU, Quantity); call sender.Send(req). [HandleResponse] map ProcessOrderItemResponse →
// OrderItemResult: copy OrderID, ItemID, AvailableQty, Reserved, Status, UnitPrice from response; push downstream via
// sc.Collect. [EndRequest] log the outcome.
type ProcessOrderItemSink struct{}

// BeginRequest is called once per gRPC stream/call, before any ConsumeMessage.
// MUST:
//   - Attach outgoing metadata (trace ID, auth) to ctx via metadata.AppendToOutgoingContext
//   - Initialize ProcessOrderItemSinkHandlerState with per-stream state
//   - Return non-nil error to abort — framework will NOT call ConsumeMessage,
//     HandleResponse, or EndRequest; release any acquired resources here.
func (ep *ProcessOrderItemSink) BeginRequest(ctx context.Context, _ datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error]) (context.Context, ProcessOrderItemSinkHandlerState, error) {
	//TODO: attach outgoing metadata to ctx, initialize ProcessOrderItemSinkHandlerState
	return ctx, ProcessOrderItemSinkHandlerState{}, nil
}

// ConsumeMessage converts value *types.OrderItem into a gRPC request and sends it.
// Map value to *processorderitem.ProcessOrderItemRequest, call sender.Send. Return non-nil error to abort.
// For bidi/client-streaming: call resultCtx.Done() on the LAST message to close the send
// side and let the framework await responses. For unary/server-streaming Done() is a no-op.
func (ep *ProcessOrderItemSink) ConsumeMessage(ctx context.Context, _ datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemSinkHandlerState, value *types.OrderItem, sender datasinkgrpc.Sender[*processorderitem.ProcessOrderItemRequest], resultCtx datasinkgrpc.ResultContext) error {
	var req *processorderitem.ProcessOrderItemRequest //TODO: map value (*types.OrderItem) → *processorderitem.ProcessOrderItemRequest
	if err := sender.Send(ctx, req); err != nil {
		return err
	}
	// Call Done() on the last message for bidi/client-streaming to close the send side.
	// For unary/server-streaming this is a no-op and may be omitted.
	resultCtx.Done()
	return nil
}

// HandleResponse is called for each response received from the gRPC server.
// For unary: called once. For server-streaming/bidi: called per response message.
// Map response to *types.OrderItemResult and push downstream via sc.Collect. Return non-nil error to abort.
//
// Thread safety: for bidi-streaming HandleResponse may be called CONCURRENTLY
// with ConsumeMessage — synchronise access to ProcessOrderItemSinkHandlerState accordingly.
func (ep *ProcessOrderItemSink) HandleResponse(ctx context.Context, sc datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemSinkHandlerState, response *processorderitem.ProcessOrderItemResponse) error {
	var result *types.OrderItemResult //TODO: map response (*processorderitem.ProcessOrderItemResponse) → *types.OrderItemResult
	sc.Collect(ctx, result)
	return nil
}

// EndRequest finalises the stream after all responses are received (or on error).
// err is the first non-nil error from any earlier stage; nil on the happy path.
// Does NOT return an error — log or record metrics here; release resources.
func (ep *ProcessOrderItemSink) EndRequest(_ context.Context, _ datasinkgrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], err error, _ ProcessOrderItemSinkHandlerState) {
	//TODO: release resources, log outcome (err != nil means call failed)
}

// MakeProcessOrderItemSink implements the handler for the ProcessOrderItem gRPC sink endpoint.
// It sends outgoing gRPC calls for each message received from the stream.
// Instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeProcessOrderItemSink(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.GrpcEndpointConfig) (*ProcessOrderItemSink, error) {
	return &ProcessOrderItemSink{}, nil
}
