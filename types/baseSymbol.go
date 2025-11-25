package types

// BaseSymbolParams defines the common symbol parameters used across
// Tabdeal REST endpoints that operate on a specific trading pair.
//
// Tabdeal uses two symbol formats:
//
//  1. `symbol`
//     - The compact, Binance-style trading symbol
//     - Example: "BTCIRT", "ETHUSDT"
//
//  2. `tabdealSymbol`
//     - The underscore-separated symbol format used internally by Tabdeal
//     - Example: "BTC_IRT", "ETH_USDT"
//
// Most Tabdeal endpoints accept *either* one.
// If both are provided, Tabdeal prioritizes `symbol`.
//
// Users of this SDK should generally provide only one:
//   - Prefer `symbol` for trade, order, and market-data endpoints.
//   - Use `tabdealSymbol` only when an endpoint explicitly documents it.
//
// Fields are tagged with `omitempty` so they are omitted when empty,
// matching Tabdeal's expectations for signed request construction.
type BaseSymbolParams struct {
	Symbol        string `json:"symbol,omitempty"`
	TabdealSymbol string `json:"tabdealSymbol,omitempty"`
}
