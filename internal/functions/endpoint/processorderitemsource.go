package endpoint

import (
	"context"

	datasourcegrpc "github.com/gorundebug/servicelib/datasource/grpc"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"

	"github.com/gorundebug/inventory_service_api/pkg/generated/proto/inventoryserviceapi/processorderitem"
	"github.com/gorundebug/model_go/pkg/types"
)

// ProcessOrderItemSourceType is the typed handler function for the ProcessOrderItemSource gRPC endpoint.
type ProcessOrderItemSourceType = datasourcegrpc.UnaryHandler[*processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse]

// processOrderItemHandler is a type alias for the EndpointHandler generic instantiation used throughout this file.
type processOrderItemHandler = datasourcegrpc.EndpointHandler[ProcessOrderItemSourceHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, error]

var _ processOrderItemHandler = (*ProcessOrderItemSource)(nil)

func MakeEndpointConsumerProcessOrderItemSource(stream runtime.TypedInputStream[*types.OrderItem, *types.OrderItemResult, error], handler processOrderItemHandler) (runtime.Consumer[*types.OrderItem], ProcessOrderItemSourceType, error) {
	return datasourcegrpc.MakeGRPCNoStreamingEndpointConsumer[ProcessOrderItemSourceHandlerState, *processorderitem.ProcessOrderItemRequest, *processorderitem.ProcessOrderItemResponse, *types.OrderItem, *types.OrderItemResult, error](stream, handler)
}

// ProcessOrderItemSourceHandlerState holds per-request state created by BeginRequest for each incoming gRPC request.
type ProcessOrderItemSourceHandlerState struct {
}

// ProcessOrderItemSource
type ProcessOrderItemSource struct{}

func (ep *ProcessOrderItemSource) BeginRequest(ctx context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error]) (context.Context, ProcessOrderItemSourceHandlerState, error) {
	return ctx, ProcessOrderItemSourceHandlerState{}, nil
}

func (ep *ProcessOrderItemSource) ConsumeMessage(ctx context.Context, sc datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], handlerState ProcessOrderItemSourceHandlerState, req *processorderitem.ProcessOrderItemRequest, resultCtx datasourcegrpc.ResultContext[ProcessOrderItemSourceHandlerState, *types.OrderItem, *processorderitem.ProcessOrderItemResponse, *types.OrderItemResult, error], sender datasourcegrpc.Sender[*types.OrderItemResult, *processorderitem.ProcessOrderItemResponse]) (context.Context, error) {
	item := &types.OrderItem{
		OrderID:  req.OrderId,
		ItemID:   req.ItemId,
		SKU:      req.Sku,
		Quantity: int(req.Quantity),
	}

	resultCtx.SetResultCallback(req.ItemId, func(ctx context.Context, sc datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemSourceHandlerState, value *types.OrderItemResult, sender datasourcegrpc.Sender[*types.OrderItemResult, *processorderitem.ProcessOrderItemResponse]) bool {
		_ = sender.Send(ctx, &processorderitem.ProcessOrderItemResponse{
			AvailableQty: int32(value.AvailableQty),
			Reserved:     value.Reserved,
			Status:       value.Status,
		})
		return true
	})

	sc.Collect(ctx, item)
	return ctx, nil
}

func (ep *ProcessOrderItemSource) GetMessageID(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemSourceHandlerState, value *types.OrderItemResult) string {
	return value.ItemID
}

func (ep *ProcessOrderItemSource) Eof(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemSourceHandlerState) {
}

func (ep *ProcessOrderItemSource) EndRequest(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], err error, _ ProcessOrderItemSourceHandlerState) error {
	return err
}

// MakeProcessOrderItemSource implements the handler for the ProcessOrderItemSource gRPC source endpoint.
func MakeProcessOrderItemSource(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.GrpcEndpointConfig) (*ProcessOrderItemSource, error) {
	return &ProcessOrderItemSource{}, nil
}
