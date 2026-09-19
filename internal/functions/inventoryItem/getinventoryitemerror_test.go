package inventoryItem

import (
	"context"
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
	value := &inventoryFailure{orderID: "order-1", itemID: "item-1", sku: "SKU-001", requestedQty: 3, availableQty: 1, unitPrice: 2.5}
	f.Map(context.Background(), nil, value, out)
	assert.Len(t, collected, 1)
	assert.Equal(t, "order-1", collected[0].OrderID)
	assert.Equal(t, "OUT_OF_STOCK", collected[0].Status)
}
