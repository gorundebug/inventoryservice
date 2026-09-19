package inventoryItem

import (
	"context"

	"github.com/gorundebug/model_go/pkg/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.MapFunction[error, *types.OrderItemResult] = (*GetInventoryItemError)(nil)

// GetInventoryItemError
type GetInventoryItemError struct{}

func (f *GetInventoryItemError) Map(ctx context.Context, _ runtime.Stream, value error, out runtime.Collect[*types.OrderItemResult]) {
	failure, ok := value.(*inventoryFailure)
	if !ok {
		out.Out(ctx, &types.OrderItemResult{Status: "PROCESSING_ERROR", Error: value.Error()})
		return
	}
	out.Out(ctx, &types.OrderItemResult{
		OrderID: failure.orderID, ItemID: failure.itemID, SKU: failure.sku,
		RequestedQty: failure.requestedQty, AvailableQty: failure.availableQty,
		Reserved: false, Status: "OUT_OF_STOCK", UnitPrice: failure.unitPrice,
		Error: failure.Error(),
	})
}

// MakeGetInventoryItemError is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeGetInventoryItemError(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.MapStreamConfig) (*GetInventoryItemError, error) {
	return &GetInventoryItemError{}, nil
}
