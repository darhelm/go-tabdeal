package types

// Wallet represents a single wallet entry for a specific asset in the
// user's Tabdeal account. Each entry reports the available balance
// ("free") as well as the blocked balance ("freeze"), which includes
// amounts locked by open orders or pending operations.
//
// All balance values are returned as strings to preserve precision.
type Wallet struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Freeze string `json:"freeze"`
}

// GetWalletParams defines the optional query parameters used when
// requesting wallet information from Tabdeal.
//
// If Asset is provided, the response is restricted to that asset.
// If Asset is omitted, Tabdeal returns the full list of wallet entries.
type GetWalletParams struct {
	Asset string `json:"asset,omitempty"`
}
