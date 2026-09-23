package inventoryItem

import (
	"context"
	inventorytypes "github.com/gorundebug/inventoryservice/internal/types"

	"github.com/gorundebug/model_go/pkg/types"

	"github.com/gorundebug/servicelib/runtime"

	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.MapFunction[*inventorytypes.InventoryFailure, *types.OrderItemResult] = (*GetInventoryItemError)(nil)

// GetInventoryItemError
type GetInventoryItemError struct{}

func (f *GetInventoryItemError) Map(ctx context.Context, _ runtime.Stream, value *inventorytypes.InventoryFailure, out runtime.Collect[*types.OrderItemResult]) {
	item := value.Item
	out.Out(ctx, &types.OrderItemResult{
		OrderID: item.OrderID, ItemID: item.ItemID, SKU: item.SKU,
		RequestedQty: item.Quantity, AvailableQty: value.AvailableQty,
		Reserved: false, Status: "OUT_OF_STOCK", UnitPrice: item.UnitPrice,
		Error: "inventory is out of stock",
	})

}

// MakeGetInventoryItemError is instantiated once at application startup via its maker function.
// Fields of this struct are not protected by any synchronization — do not use
// shared mutable state here without external synchronization.
func MakeGetInventoryItemError(ctx context.Context, env environment.ServiceEnvironment) (*GetInventoryItemError, error) {
	return &GetInventoryItemError{}, nil
}
