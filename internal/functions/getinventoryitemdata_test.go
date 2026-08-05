package functions

import (
	"context"
	"testing"

	"github.com/gorundebug/model/pkg/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Look up the inventory record by OrderItem.SKU. If stock >= OrderItem.Quantity: reserve the stock atomically and emit
// OrderItemResult{Reserved: true, Status: CONFIRMED, AvailableQty: OrderItem.Quantity} via out. If stock is
// insufficient: emit OrderItemResult{Reserved: false, Status: OUT_OF_STOCK, AvailableQty: actual available} via the
// error output stream (rout) so the merge node can still count it as a processed item.

func makeTestInventory() *GetInventoryItemData {
	return &GetInventoryItemData{
		stock: map[string]int{
			"SKU-001": 10,
			"SKU-002": 0,
		},
	}
}

func TestGetInventoryItemData_Process_Confirmed(t *testing.T) {
	f := makeTestInventory()
	var results []*types.OrderItemResult
	out := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		results = append(results, v)
	})
	var errs []*types.OrderItemResult
	rout := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		errs = append(errs, v)
	})
	value := &types.OrderItem{OrderID: "order-1", ItemID: "item-1", SKU: "SKU-001", Quantity: 3}
	f.Process(context.Background(), nil, value, out, rout)
	assert.Empty(t, errs)
	assert.Len(t, results, 1)
	assert.True(t, results[0].Reserved)
	assert.Equal(t, "CONFIRMED", results[0].Status)
	assert.Equal(t, 3, results[0].AvailableQty)
}

func TestGetInventoryItemData_Process_OutOfStock(t *testing.T) {
	f := makeTestInventory()
	var results []*types.OrderItemResult
	out := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		results = append(results, v)
	})
	var errs []*types.OrderItemResult
	rout := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		errs = append(errs, v)
	})
	value := &types.OrderItem{OrderID: "order-2", ItemID: "item-2", SKU: "SKU-002", Quantity: 1}
	f.Process(context.Background(), nil, value, out, rout)
	assert.Empty(t, results)
	assert.Len(t, errs, 1)
	assert.False(t, errs[0].Reserved)
	assert.Equal(t, "OUT_OF_STOCK", errs[0].Status)
	assert.Equal(t, 0, errs[0].AvailableQty)
}

func TestGetInventoryItemData_Process_InsufficientStock(t *testing.T) {
	f := makeTestInventory()
	var results []*types.OrderItemResult
	out := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		results = append(results, v)
	})
	var errs []*types.OrderItemResult
	rout := runtime.CollectFunc[*types.OrderItemResult](func(_ context.Context, v *types.OrderItemResult) {
		errs = append(errs, v)
	})
	value := &types.OrderItem{OrderID: "order-3", ItemID: "item-3", SKU: "SKU-001", Quantity: 100}
	f.Process(context.Background(), nil, value, out, rout)
	assert.Empty(t, results)
	assert.Len(t, errs, 1)
	assert.False(t, errs[0].Reserved)
	assert.Equal(t, "OUT_OF_STOCK", errs[0].Status)
}
