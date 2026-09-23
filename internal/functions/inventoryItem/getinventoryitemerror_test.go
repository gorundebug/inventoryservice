package inventoryItem

import (
	"context"
	inventorytypes "github.com/gorundebug/inventoryservice/internal/types"
	"testing"

	"github.com/gorundebug/model_go/pkg/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// When inventory processing fails, return an OUT_OF_STOCK result with no available quantity. Preserve the order and
// item identity and requested quantity, and record the failure.

func TestGetInventoryItemError_Map(t *testing.T) {
	f := &GetInventoryItemError{}
	var collected []*types.OrderItemResult
	out := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		collected = append(collected, v)
	})
	value := &inventorytypes.InventoryFailure{Item: &types.OrderItem{OrderID: "order-1", ItemID: "item-1", SKU: "SKU-001", Quantity: 3, UnitPrice: 2.5}, AvailableQty: 1}
	f.Map(context.Background(), nil, value, out)
	assert.Equal(t, []*types.OrderItemResult{{
		OrderID: "order-1", ItemID: "item-1", SKU: "SKU-001",
		RequestedQty: 3, AvailableQty: 1, UnitPrice: 2.5,
		Reserved: false, Status: "OUT_OF_STOCK", Error: "inventory is out of stock",
	}}, collected)
}
