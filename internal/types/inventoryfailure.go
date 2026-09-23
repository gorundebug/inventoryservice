package types

import modeltypes "github.com/gorundebug/model_go/pkg/types"

// InventoryFailure is a business outcome, not a transport error.
type InventoryFailure struct {
	Item         *modeltypes.OrderItem
	AvailableQty int
}
