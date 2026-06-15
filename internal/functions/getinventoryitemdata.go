package functions

import (
	"context"

	"github.com/gorundebug/model/pkg/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.ProcessFunction[*types.OrderItem, *types.OrderItemResult, *types.OrderItemResult] = (*GetInventoryItemData)(nil)

// GetInventoryItemData
type GetInventoryItemData struct{}

func (f *GetInventoryItemData) Process(_ context.Context, _ runtime.Stream, value *types.OrderItem, out runtime.Collect[*types.OrderItemResult], rout runtime.Collect[*types.OrderItemResult]) {
	//TODO: Need to be implemented
	// Look up the inventory record by OrderItem.SKU; retrieve current stock and UnitPrice from the record. Always copy
	// OrderID, ItemID, SKU, RequestedQty (=OrderItem.Quantity), UnitPrice into the result. If stock >= OrderItem.Quantity:
	// reserve the stock atomically and emit OrderItemResult{OrderID, ItemID, SKU, RequestedQty, UnitPrice, Reserved: true,
	// Status: CONFIRMED, AvailableQty: OrderItem.Quantity} via out. If stock is insufficient: emit OrderItemResult{OrderID,
	// ItemID, SKU, RequestedQty, UnitPrice, Reserved: false, Status: OUT_OF_STOCK, AvailableQty: actual available} via
	// rout.
}

// MakeGetInventoryItemData is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeGetInventoryItemData(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.ProcessStreamConfig) (*GetInventoryItemData, error) {
	return &GetInventoryItemData{}, nil
}
