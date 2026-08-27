package inventoryItem

import (
	"context"
	"sync/atomic"

	"github.com/gorundebug/model/pkg/types"

	"github.com/gorundebug/servicelib/runtime"
	runtimecfg "github.com/gorundebug/servicelib/runtime/config"
	"github.com/gorundebug/servicelib/runtime/environment"
	"github.com/gorundebug/servicelib/transformation"
)

var _ transformation.ProcessFunction[*types.OrderItem, *types.OrderItemResult, *types.OrderItemResult] = (*GetInventoryItemData)(nil)

// GetInventoryItemData
type GetInventoryItemData struct {
	stock map[string]*atomic.Int64 // immutable SKU index; atomic quantity per SKU
}

func (f *GetInventoryItemData) Process(ctx context.Context, _ runtime.Stream, value *types.OrderItem, out runtime.Collect[*types.OrderItemResult], rout runtime.Collect[*types.OrderItemResult]) {
	stock, ok := f.stock[value.SKU]
	if ok {
		quantity := int64(value.Quantity)
		for available := stock.Load(); available >= quantity; available = stock.Load() {
			if stock.CompareAndSwap(available, available-quantity) {
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
		}
	}
	available := int64(0)
	if ok {
		available = stock.Load()
	}
	rout.Out(ctx, &types.OrderItemResult{
		OrderID:      value.OrderID,
		ItemID:       value.ItemID,
		SKU:          value.SKU,
		RequestedQty: value.Quantity,
		AvailableQty: int(available),
		Reserved:     false,
		Status:       "OUT_OF_STOCK",
		UnitPrice:    value.UnitPrice,
	})
}

// MakeGetInventoryItemData is instantiated once at application startup via its maker function.
func MakeGetInventoryItemData(ctx context.Context, env environment.ServiceEnvironment, cfg *runtimecfg.ProcessStreamConfig) (*GetInventoryItemData, error) {
	return &GetInventoryItemData{
		stock: map[string]*atomic.Int64{
			"SKU-001": atomicInt64(100),
			"SKU-002": atomicInt64(50),
			"SKU-003": atomicInt64(25),
		},
	}, nil
}

func atomicInt64(value int64) *atomic.Int64 {
	result := &atomic.Int64{}
	result.Store(value)
	return result
}
