package functions

import (
	"context"
	"sync"

	"github.com/gorundebug/model/pkg/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.ProcessFunction[*types.OrderItem, *types.OrderItemResult, *types.OrderItemResult] = (*GetInventoryItemData)(nil)

// GetInventoryItemData
type GetInventoryItemData struct {
	mu    sync.Mutex
	stock map[string]int // SKU -> available quantity
}

func (f *GetInventoryItemData) Process(ctx context.Context, _ runtime.Stream, value *types.OrderItem, out runtime.Collect[*types.OrderItemResult], rout runtime.Collect[*types.OrderItemResult]) {
	f.mu.Lock()
	available, ok := f.stock[value.SKU]
	if ok && available >= value.Quantity {
		f.stock[value.SKU] -= value.Quantity
		f.mu.Unlock()
		out.Out(ctx, &types.OrderItemResult{
			OrderID:      value.OrderID,
			ItemID:       value.ItemID,
			SKU:          value.SKU,
			RequestedQty: value.Quantity,
			AvailableQty: value.Quantity,
			Reserved:     true,
			Status:       "CONFIRMED",
			UnitPrice:    value.UnitPrice,
		})
		return
	}
	avail := available
	if !ok {
		avail = 0
	}
	f.mu.Unlock()
	rout.Out(ctx, &types.OrderItemResult{
		OrderID:      value.OrderID,
		ItemID:       value.ItemID,
		SKU:          value.SKU,
		RequestedQty: value.Quantity,
		AvailableQty: avail,
		Reserved:     false,
		Status:       "OUT_OF_STOCK",
		UnitPrice:    value.UnitPrice,
	})
}

// MakeGetInventoryItemData is instantiated once at application startup via its maker function.
func MakeGetInventoryItemData(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.ProcessStreamConfig) (*GetInventoryItemData, error) {
	return &GetInventoryItemData{
		stock: map[string]int{
			"SKU-001": 100,
			"SKU-002": 50,
			"SKU-003": 25,
		},
	}, nil
}
