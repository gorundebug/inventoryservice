package functions

import (
	"context"
	"testing"

	"github.com/gorundebug/model/pkg/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Look up the inventory record by OrderItem.SKU; retrieve current stock and UnitPrice from the record. Always copy
// OrderID, ItemID, SKU, RequestedQty (=OrderItem.Quantity), UnitPrice into the result. If stock >= OrderItem.Quantity:
// reserve the stock atomically and emit OrderItemResult{OrderID, ItemID, SKU, RequestedQty, UnitPrice, Reserved: true,
// Status: CONFIRMED, AvailableQty: OrderItem.Quantity} via out. If stock is insufficient: emit OrderItemResult{OrderID,
// ItemID, SKU, RequestedQty, UnitPrice, Reserved: false, Status: OUT_OF_STOCK, AvailableQty: actual available} via
// rout.

func TestGetInventoryItemData_Process(t *testing.T) {
	t.Skip("not yet implemented") // TODO: remove when implementation is ready
	f := &GetInventoryItemData{}
	var results []*types.OrderItemResult
	out := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		results = append(results, v)
	})
	var errs []*types.OrderItemResult
	rout := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		errs = append(errs, v)
	})
	var value *types.OrderItem
	// TODO: populate value with meaningful test data
	f.Process(context.Background(), nil, value, out, rout)
	assert.Empty(t, errs)
	assert.NotEmpty(t, results)
}
