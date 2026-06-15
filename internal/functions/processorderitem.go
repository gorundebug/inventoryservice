package functions

import (
	"context"

	datasourcegrpc "github.com/gorundebug/servicelib/datasource/grpc"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/runtime/environment/log"

	"github.com/gorundebug/inventory_service_api/pkg/generated/proto/inventoryserviceapi/processorderitem"
	"github.com/gorundebug/model/pkg/types"
)

// ProcessOrderItemType is the typed handler function for the ProcessOrderItem gRPC endpoint.
type ProcessOrderItemType = datasourcegrpc.UnaryHandler[*processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse]

// processOrderItemHandler is a type alias for the EndpointHandler generic instantiation used throughout this file.
type processOrderItemHandler = datasourcegrpc.EndpointHandler[ProcessOrderItemHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, error]

var _ processOrderItemHandler = (*ProcessOrderItem)(nil)

func MakeEndpointConsumerProcessOrderItem(stream runtime.TypedInputStream[*types.OrderItem, *types.OrderItemResult, error], handler processOrderItemHandler) (runtime.Consumer[*types.OrderItem], ProcessOrderItemType, error) {
	return datasourcegrpc.MakeGRPCNoStreamingEndpointConsumer[ProcessOrderItemHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, error](stream, handler)
}

// ProcessOrderItemHandlerState holds per-request state created by BeginRequest for each incoming gRPC request.
// Enables safe concurrent processing — no synchronization needed between requests.
// Add fields here to carry data across BeginRequest → ConsumeMessage → Eof → EndRequest.
type ProcessOrderItemHandlerState struct {
}

// ProcessOrderItem
// Outgoing unary gRPC call to the Inventory Service. [ConsumeMessage] map OrderItem → ProcessOrderItemRequest
// (OrderID, ItemID, SKU, Quantity); call sender.Send(req). [HandleResponse] map ProcessOrderItemResponse →
// OrderItemResult: copy OrderID, ItemID, AvailableQty, Reserved, Status, UnitPrice from response; push downstream via
// sc.Collect. [EndRequest] log the outcome.
type ProcessOrderItem struct{}

// BeginRequest is called once per RPC, before any ConsumeMessage calls.
// It is the first line of defence: use it to enforce routing rules, apply rate limiting,
// check auth tokens, or extract trace IDs from gRPC metadata — anything that should
// gate the RPC before it enters the pipeline.
// Store whatever ConsumeMessage will need (e.g. auth context, parsed metadata) in ProcessOrderItemHandlerState.
// Return non-nil error to reject the RPC immediately with a gRPC status error.
//
// IMPORTANT: if a non-nil error is returned, the framework does NOT call
// ConsumeMessage, Eof, or EndRequest. Release any acquired resources here.
func (ep *ProcessOrderItem) BeginRequest(ctx context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error]) (context.Context, ProcessOrderItemHandlerState, error) {
	//TODO: extract metadata from ctx, initialize ProcessOrderItemHandlerState
	// Optional: if the downstream link uses a priority task pool and you need
	// per-request priority (e.g. based on customer tier or SLA level), wrap ctx:
	//   ctx = runtime.WithPriority(ctx, priority)
	// Otherwise the pool uses the priority configured in the link settings.
	return ctx, ProcessOrderItemHandlerState{}, nil
}

// ConsumeMessage is called once per inbound gRPC message.
// Pattern: map *processorderitem.ProcessOrderItemRequest → *types.OrderItem, emit via sc.Collect, register SetResultCallback.
// If a non-nil error is returned, Eof will NOT be called; EndRequest is called with that error.
// You MAY still write to sender before returning the error.
//
// Done() semantics:
//   - Unary: no-op; the framework waits for a reply written to sender.
//   - Client-streaming: closes the send side (no further responses); framework sends a zero
//     response via SendAndClose if sender.Send was never called.
//   - Server-streaming / bidi-streaming: stops the pipeline (no further results expected).
func (ep *ProcessOrderItem) ConsumeMessage(ctx context.Context, sc datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState, req *processorderitem.ProcessOrderItemRequest, resultCtx datasourcegrpc.ResultContext[ProcessOrderItemHandlerState, *types.OrderItem, *processorderitem.ProcessOrderItemResponse, *types.OrderItemResult, error], sender datasourcegrpc.Sender[*types.OrderItemResult, *processorderitem.ProcessOrderItemResponse]) (context.Context, error) {
	//TODO: map req (*processorderitem.ProcessOrderItemRequest) → *types.OrderItem
	var value *types.OrderItem
	id := "" //TODO: set to value that matches ep.GetMessageID for the corresponding *types.OrderItemResult result
	resultCtx.SetResultCallback(id, func(ctx context.Context, sc datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState, result *types.OrderItemResult, sender datasourcegrpc.Sender[*types.OrderItemResult, *processorderitem.ProcessOrderItemResponse]) bool {
		var resp *processorderitem.ProcessOrderItemResponse //TODO: build response from result
		if err := sender.Send(ctx, resp); err != nil {
			sc.Logger.Error(ctx, "failed to send response", log.Err(err))
			resultCtx.Done()
			return true // deregister callback — no more results expected
		}
		// Done signals that no further results are expected for this request.
		// For unary: no-op (sender.Send already unblocks the framework).
		// For client-streaming: closes the send side.
		// For server-streaming/bidi: stops the pipeline.
		// Call it after sending the last response to the client.
		resultCtx.Done()
		// Return true to deregister this callback — no more results are expected.
		// Return false if the pipeline may emit additional results for this request.
		return true
	})
	sc.Collect(ctx, value)
	return ctx, nil
}

// GetMessageID returns a stable correlation ID from *types.OrderItemResult so the framework can
// route async pipeline results back to the correct ResultContext callback.
// Return "" if this endpoint does not use async result routing.
func (ep *ProcessOrderItem) GetMessageID(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState, _ *types.OrderItemResult) string {
	//TODO: return stable ID from *types.OrderItemResult (e.g. request_id field), or ""
	return ""
}

// Eof signals the end of the inbound message stream. Called for ALL RPC types:
//   - Unary / server-streaming: called once, immediately after the single ConsumeMessage
//   - Client-streaming / bidi: called when the client closes its send side (last Recv → io.EOF)
//
// NOT called if ConsumeMessage returned an error.
// Use it to finalise aggregation or flush buffered state before EndRequest.
func (ep *ProcessOrderItem) Eof(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState) {
	//TODO: finalise aggregation or flush buffered state before EndRequest
}

// EndRequest finalises the RPC. Called in all exit paths: after Eof, after a Recv or
// ConsumeMessage error, or after context cancellation.
// Its return value is propagated as the final gRPC status to the client.
//
//   - err == nil: happy path — response was already sent via sender in SetResultCallback.
//     Use EndRequest for cleanup only. Do NOT return a non-nil error here: it would
//     override the successful response already delivered to the client.
//   - err != nil: a stage failed or the context was cancelled.
//     Return a gRPC status error to signal failure to the client.
func (ep *ProcessOrderItem) EndRequest(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], err error, _ ProcessOrderItemHandlerState) error {
	if err != nil {
		return err //TODO: optionally wrap as a gRPC status error (e.g. status.Errorf)
	}
	return nil
}

// MakeProcessOrderItem implements the handler for the ProcessOrderItem gRPC source endpoint.
// It processes incoming gRPC requests and produces messages into the stream.
// Instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeProcessOrderItem(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.GrpcEndpointConfig) (*ProcessOrderItem, error) {
	return &ProcessOrderItem{}, nil
}
