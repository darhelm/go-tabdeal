package types

// BaseOrderResponse represents the common fields returned by Tabdeal
// for order-related endpoints, including newly-created orders, queried
// orders, and updated order states.
//
// Numeric and monetary values are returned as strings to preserve
// precision. Time fields represent milliseconds since Unix epoch.
type BaseOrderResponse struct {
	Symbol               string `json:"symbol"`
	TabdealSymbol        string `json:"tabdealSymbol"`
	OrderId              int64  `json:"orderId"`
	OrderListId          int64  `json:"orderListId"`
	ClientOrderId        string `json:"clientOrderId,omitempty"`
	TransactTime         int64  `json:"transactTime"`
	Price                string `json:"price"`
	OrigQty              string `json:"origQty"`
	ExecutedQty          string `json:"executedQty"`
	CummulativeQuoteQty  string `json:"cummulativeQuoteQty"`
	CumulativeQuoteQty   string `json:"cumulativeQuoteQty"`
	Status               string `json:"status"`
	Type                 string `json:"type"`
	Side                 string `json:"side"`
	StopPrice            string `json:"stopPrice"`
	UpdateTime           int64  `json:"updateTime"`
	IsWorking            bool   `json:"isWorking"`
	IsStopOrderTriggered bool   `json:"isStopOrderTriggered"`
}

// Fills represents an individual trade execution that occurred while
// fulfilling an order. A single order may generate multiple fills,
// each contributing to the total executed quantity.
type Fills struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	TradeId         int64  `json:"tradeId"`
}

// CreateOrderResponse contains the complete details of an order
// immediately after it is submitted to Tabdeal. This includes the
// initial order fields as well as any execution fills generated at
// the time of placement.
type CreateOrderResponse struct {
	BaseOrderResponse
	Fills []Fills `json:"fills"`
}

// CreateOrderParams defines the parameters required to create a new
// order on Tabdeal. The meaning of fields depends on the order type:
//
//   - LIMIT orders require price and quantity.
//   - MARKET orders require quantity only.
//   - STOP or STOP-LIMIT orders may require stopPrice.
//
// Side and type must correspond to the allowed values returned by
// Tabdeal's market-information endpoint.
//
// newClientOrderId may be supplied to assign a custom tracking ID
// to the order.
type CreateOrderParams struct {
	BaseSymbolParams

	Side             string  `json:"side"`
	Type             string  `json:"type"`
	Quantity         float64 `json:"quantity"`
	NewClientOrderId string  `json:"newClientOrderId,omitempty"`
	Price            float64 `json:"price,omitempty"`
	StopPrice        float64 `json:"stopPrice,omitempty"`
}

// GetOrderStatusParams specifies how to retrieve the status of a single
// order. The order may be identified by either:
//
//   - orderId
//   - origClientOrderId
//
// At least one should be supplied.
type GetOrderStatusParams struct {
	OrderId           int    `json:"orderId,omitempty"`
	OrigClientOrderId string `json:"origClientOrderId,omitempty"`
}

// OrderStatusResponse provides detailed information about a specific
// order, as returned by Tabdeal's order-status query endpoint.
// Includes order fields plus the fee associated with completed orders.
type OrderStatusResponse struct {
	BaseOrderResponse
	Fee string `json:"fee"`
}

// GetOpenOrdersParams defines the optional parameters used when
// requesting the list of currently open orders. When a symbol is
// provided, only orders in that market are returned.
type GetOpenOrdersParams struct {
	BaseSymbolParams
}

// UserTradeResponse represents a single user-specific trade execution,
// including price, quantity, commission, and maker/taker flag.
// These records are returned by Tabdeal's user-trade history endpoints.
type UserTradeResponse struct {
	Symbol          string `json:"symbol"`
	TabdealSymbol   string `json:"tabdealSymbol"`
	Id              int64  `json:"id"`
	OrderId         int64  `json:"orderId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
}

// GetUserOrdersHistoryParams defines filters for querying a user's
// historical order activity. The range may be limited using startTime
// and endTime. Limit controls how many records are returned.
type GetUserOrdersHistoryParams struct {
	BaseSymbolParams
	StartTime int64 `json:"startTime,omitempty"`
	EndTime   int64 `json:"endTime,omitempty"`
	Limit     int64 `json:"limit,omitempty"`
}

// GetUserTradesParams extends historical-order filters with the ability
// to return trades associated with a specific orderId, enabling finer
// selection when analyzing past executions.
type GetUserTradesParams struct {
	GetUserOrdersHistoryParams
	OrderId int64 `json:"orderId,omitempty"`
}
