package types

// Filter represents a ruleset applied by Tabdeal to a given trading pair.
// Each filter describes a specific constraint related to price, quantity,
// order size, or other market-level validation.
//
// The exact fields populated depend on the filterType. Examples include:
//
//	PRICE_FILTER:
//	  - minPrice
//	  - maxPrice
//	  - tickSize
//
//	PERCENT_PRICE:
//	  - multiplierUp
//	  - multiplierDown
//	  - avgPriceMins
//
//	LOT_SIZE / MARKET_LOT_SIZE:
//	  - minQty
//	  - maxQty
//	  - stepSize
//
//	MIN_NOTIONAL:
//	  - minNotional
//	  - applyToMarket
//
// Some filters return numeric values as strings. Optional fields use pointers
// where necessary to distinguish between zero values and absent fields.
type Filter struct {
	FilterType string `json:"filterType"`

	// PRICE_FILTER
	MinPrice string `json:"minPrice,omitempty"`
	MaxPrice string `json:"maxPrice,omitempty"`
	TickSize string `json:"tickSize,omitempty"`

	// PERCENT_PRICE
	MultiplierUp   float64 `json:"multiplierUp,omitempty"`
	MultiplierDown float64 `json:"multiplierDown,omitempty"`
	AvgPriceMins   int64   `json:"avgPriceMins,omitempty"`

	// LOT_SIZE and MARKET_LOT_SIZE
	MinQty   string `json:"minQty,omitempty"`
	MaxQty   string `json:"maxQty,omitempty"`
	StepSize string `json:"stepSize,omitempty"`

	// MIN_NOTIONAL
	MinNotional   string `json:"minNotional,omitempty"`
	ApplyToMarket bool   `json:"applyToMarket,omitempty"`
}

// MarketInformation describes a single market available on Tabdeal,
// including symbol information, trading permissions, supported order types,
// and the complete set of validation filters associated with the market.
//
// This structure maps directly to the result returned by Tabdeal's
// market-information endpoint.
type MarketInformation struct {
	Symbol                     string   `json:"symbol"`
	TabdealSymbol              string   `json:"tabdealSymbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         string   `json:"baseAssetPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuoteAssetPrecision        string   `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    string   `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   string   `json:"quoteCommissionPrecision"`
	OrderTypes                 []string `json:"orderTypes"`
	IcebergAllowed             bool     `json:"icebergAllowed"`
	OcoAllowed                 bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop          bool     `json:"allowTrailingStop"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	Filters                    []Filter `json:"filters"`
	Permissions                []string `json:"permissions"`
}

// OrderBook represents the current aggregated order book for a market.
// Each entry in Asks and Bids is a [price, quantity] pair formatted as strings.
//
// The best ask appears at Asks[0], and the best bid appears at Bids[0].
type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

// Trade describes a single executed trade on Tabdeal's spot market.
// Tabdeal returns price, quantity, and quote quantity as strings,
// along with timestamp and taker/maker information.
type Trade struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
}

// ServerTime represents the server's current timestamp in milliseconds,
// as returned by Tabdeal's time-synchronization endpoint.
type ServerTime struct {
	ServerTime int64 `json:"serverTime"`
}

// GetRecentTradesParams defines the query parameters used when retrieving
// recent trades for a specific market. The limit parameter controls how many
// trades are returned, up to the maximum supported by Tabdeal.
type GetRecentTradesParams struct {
	BaseSymbolParams
	Limit int64 `json:"limit,omitempty"`
}

// GetOrderBookParams defines the query parameters used to fetch
// the current order book for a specific market. The limit parameter controls
// how many price levels are included in the response.
type GetOrderBookParams struct {
	BaseSymbolParams
	Limit int64 `json:"limit,omitempty"`
}
