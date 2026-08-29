package inventoryItem

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/gorundebug/model_go/pkg/types"
	"github.com/gorundebug/servicelib/runtime"
	"github.com/stretchr/testify/assert"
)

// Inventory reservations must not overdraw stock. Accepted reservations return the requested quantity;
// rejected reservations return the current available quantity through the alternate result.

func makeTestInventory() *GetInventoryItemData {
	return &GetInventoryItemData{
		stock: map[string]*atomic.Int64{
			"SKU-001": atomicInt64(10),
			"SKU-002": atomicInt64(0),
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
