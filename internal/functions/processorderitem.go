package functions

import (
	"context"

	datasourcegrpc "github.com/gorundebug/servicelib/datasource/grpc"
	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"

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
type ProcessOrderItemHandlerState struct {
}

// ProcessOrderItem
type ProcessOrderItem struct{}

func (ep *ProcessOrderItem) BeginRequest(ctx context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error]) (context.Context, ProcessOrderItemHandlerState, error) {
	return ctx, ProcessOrderItemHandlerState{}, nil
}

func (ep *ProcessOrderItem) ConsumeMessage(ctx context.Context, sc datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], handlerState ProcessOrderItemHandlerState, req *processorderitem.ProcessOrderItemRequest, resultCtx datasourcegrpc.ResultContext[ProcessOrderItemHandlerState, *types.OrderItem, *processorderitem.ProcessOrderItemResponse, *types.OrderItemResult, error], sender datasourcegrpc.Sender[*types.OrderItemResult, *processorderitem.ProcessOrderItemResponse]) (context.Context, error) {
	item := &types.OrderItem{
		OrderID:  req.OrderId,
		ItemID:   req.ItemId,
		SKU:      req.Sku,
		Quantity: int(req.Quantity),
	}

	resultCtx.SetResultCallback(req.ItemId, func(ctx context.Context, sc datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState, value *types.OrderItemResult, sender datasourcegrpc.Sender[*types.OrderItemResult, *processorderitem.ProcessOrderItemResponse]) bool {
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

func (ep *ProcessOrderItem) GetMessageID(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState, value *types.OrderItemResult) string {
	return value.ItemID
}

func (ep *ProcessOrderItem) Eof(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], _ ProcessOrderItemHandlerState) {
}

func (ep *ProcessOrderItem) EndRequest(_ context.Context, _ datasourcegrpc.StreamContext[*types.OrderItem, *types.OrderItemResult, error], err error, _ ProcessOrderItemHandlerState) error {
	return err
}

// MakeProcessOrderItem implements the handler for the ProcessOrderItem gRPC source endpoint.
func MakeProcessOrderItem(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.GrpcEndpointConfig) (*ProcessOrderItem, error) {
	return &ProcessOrderItem{}, nil
}
