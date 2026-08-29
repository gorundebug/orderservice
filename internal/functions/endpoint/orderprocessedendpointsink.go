package endpoint

import (
	"context"
	"encoding/json"

	"github.com/gorundebug/model_go/pkg/types"
	datasinkkafka "github.com/gorundebug/servicelib/datasink/kafka"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
)

// orderProcessedEndpointSinkHandler is a type alias for the EndpointHandler generic instantiation used throughout this file.
type orderProcessedEndpointSinkHandler = datasinkkafka.EndpointHandler[OrderProcessedEndpointSinkHandlerState, *types.OrderProcessed, error]

var _ orderProcessedEndpointSinkHandler = (*OrderProcessedEndpointSink)(nil)

func MakeEndpointConsumerOrderProcessedEndpointSink(stream runtime.TypedSinkStream[*types.OrderProcessed, error], handler orderProcessedEndpointSinkHandler) (runtime.Consumer[*types.OrderProcessed], error) {
	return datasinkkafka.MakeSaramaKafkaEndpointConsumer[OrderProcessedEndpointSinkHandlerState, *types.OrderProcessed, error](stream, handler)
}

// OrderProcessedEndpointSinkHandlerState holds per-stream state created by BeginRequest for each logical Kafka stream.
// Enables safe concurrent processing — no synchronization needed between streams.
// Add fields here to carry data across BeginRequest → ConsumeMessage → EndRequest.
type OrderProcessedEndpointSinkHandlerState struct {
}

// OrderProcessedEndpointSink
type OrderProcessedEndpointSink struct{}

// GetStreamID groups messages into logical streams (one BeginRequest/EndRequest per stream ID).
// Messages with the same ID share a OrderProcessedEndpointSinkHandlerState instance; return "" to route all messages to one stream.
func (ep *OrderProcessedEndpointSink) GetStreamID(_ context.Context, value *types.OrderProcessed) string {
	return value.OrderID
}

// BeginRequest is called once per stream (per unique GetStreamID), before any ConsumeMessage.
// Does NOT return an error — initialise OrderProcessedEndpointSinkHandlerState and attach outgoing metadata to ctx if needed.
// Publish the OrderProcessed event produced by MapToOrderProcessed, keyed by order ID.
func (ep *OrderProcessedEndpointSink) BeginRequest(ctx context.Context, _ runtime.Stream) (context.Context, OrderProcessedEndpointSinkHandlerState) {
	return ctx, OrderProcessedEndpointSinkHandlerState{}
}

// ConsumeMessage sends value *types.OrderProcessed as a Kafka message.
// MUST set msg.Key and/or msg.Value (both are []byte), then choose one send method:
//   - msg.Send(ctx, onDelivery) — async; onDelivery(partition, offset, err) converts delivery to error
//   - msg.SendSync(ctx)         — blocks; returns (partition, offset, error); then msg.Out(ctx, r)
//   - msg.Skip(ctx, result)     — push error downstream without sending to Kafka
//
// Return non-nil error to abort; EndRequest is called with that error.
func (ep *OrderProcessedEndpointSink) ConsumeMessage(ctx context.Context, _ runtime.Stream, _ OrderProcessedEndpointSinkHandlerState, value *types.OrderProcessed, msg *datasinkkafka.SinkMessage[error]) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg.Key = []byte(value.OrderID)
	msg.Value = payload
	msg.Send(ctx, func(_ int32, _ int64, deliveryErr error) error { return deliveryErr })
	return nil
}

// EndRequest finalises the stream after all messages are sent (or on error).
// err is the first non-nil error from ConsumeMessage; nil on the happy path.
// Does NOT return an error — log or record metrics here; release resources.
func (ep *OrderProcessedEndpointSink) EndRequest(_ context.Context, _ runtime.Stream, err error, _ OrderProcessedEndpointSinkHandlerState) {
}

// MakeOrderProcessedEndpointSink implements the handler for the PublishOrderProcessed Kafka sink endpoint.
// It produces messages to a Kafka topic for each message received from the stream.
// Instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeOrderProcessedEndpointSink(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.KafkaEndpointConfig) (*OrderProcessedEndpointSink, error) {
	return &OrderProcessedEndpointSink{}, nil
}
