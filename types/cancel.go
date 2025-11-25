package types

// CancelOrderParams defines the parameters required to cancel a single
// active order on Tabdeal's spot trading platform.
//
// An order may be identified using either of the following fields:
//
//  1. orderId
//     - The numeric order identifier assigned by Tabdeal.
//     - Example: 1234567890123
//
//  2. origClientOrderId
//     - A client-defined identifier specified during order placement.
//     - Example: "order-001"
//
// At least one identifier must be provided. If both are present,
// Tabdeal prioritizes orderId.
//
// The embedded BaseSymbolParams allows specifying the trading pair
// using either "symbol" (e.g., "BTCIRT") or "tabdealSymbol"
// (e.g., "BTC_IRT"). Only one format is required.
type CancelOrderParams struct {
	BaseSymbolParams
	OrderId           int64  `json:"orderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
}

// CancelOrderBulkParams defines the optional filtering criteria used to
// cancel multiple active orders at once.
//
// Tabdeal's bulk-cancel endpoint accepts an optional trading symbol.
// When a symbol is provided, all open orders for that market are
// cancelled. When omitted, the endpoint cancels all eligible open orders
// associated with the authenticated account.
type CancelOrderBulkParams struct {
	BaseSymbolParams
}

// CancelOrderResponse represents the standard response returned by
// Tabdeal for both single-order and bulk-order cancellation requests.
//
// Cancel endpoints will additionally include order-level fields. Such
// fields are captured through BaseOrderResponse when available.
type CancelOrderResponse struct {
	BaseOrderResponse
}
