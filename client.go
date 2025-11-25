package tabdeal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	t "github.com/darhelm/go-tabdeal/types"
	u "github.com/darhelm/go-tabdeal/utils"
)

// Constants defining the API base URL and version.
const (
	// BaseUrl is the root URL for the Tabdeal Market API.
	BaseUrl = "https://api1.tabdeal.org"
	Version = "v1"
)

// ClientOptions represents the configuration options for creating a new API client.
// These options allow customization of the HTTP client, authentication tokens
type ClientOptions struct {
	// HttpClient is the custom HTTP client to be used for API requests.
	// If nil, the default HTTP client is used.
	HttpClient *http.Client

	// Timeout specifies the request timeout duration for the HTTP client.
	Timeout time.Duration

	// BaseUrl is the base URL of the API. Defaults to the constant BaseUrl
	// if not provided.
	BaseUrl string

	// ApiKey is the token used for authenticated API requests.
	ApiKey string

	// ApiSecret is the token used for authenticated API requests.
	ApiSecret string
}

// Client represents the API client for interacting with the Tabdeal Market API.
// It manages authentication, base URL, and API requests.
type Client struct {
	// HttpClient is the HTTP client used for API requests.
	// Defaults to the Go standard library's http.DefaultClient.
	HttpClient *http.Client

	// BaseUrl is the base URL of the API used by this client.
	// Defaults to the constant BaseUrl.
	BaseUrl string

	// Version is the Api's current version
	Version string

	// ApiKey is the API key for authentication.
	ApiKey string

	// ApiSecret is the API Secret for authentication.
	ApiSecret string

	// AutoAuth enables automatic authentication if no valid tokens are provided.
	AutoAuth bool

	// AutoRefresh enables automatic refreshing of the access token when it expires.
	AutoRefresh bool
}

// NewClient initializes a new Tabdeal API client using the provided configuration
// options. It sets up the HTTP client, assigns API credentials, and applies any
// optional overrides such as custom base URLs or timeouts.
//
// Parameters:
//   - opts: ClientOptions with the following fields:
//   - HttpClient: optional custom HTTP client. If nil, a new http.Client
//     is created using opts.Timeout.
//   - Timeout: request timeout used when creating a default HttpClient.
//   - BaseUrl: optional override for the API base URL. Defaults to BaseUrl
//     ("https://api1.tabdeal.org") if empty.
//   - ApiKey: API key used for authenticated endpoints.
//   - ApiSecret: API secret used for request signing.
//
// Returns:
//   - A pointer to an initialized Client.
//   - Never returns an error, because no network operations or authentication
//     are performed inside NewClient.
//
// Behavior:
//   - If opts.BaseUrl is provided, it overrides the default BaseUrl.
//   - If opts.HttpClient is nil, a new http.Client is constructed using
//     opts.Timeout.
//   - ApiKey and ApiSecret are stored on the client for use in authenticated
//     requests.
//   - No authentication request is performed.
//   - No token refresh logic is invoked.
//
// Example:
//
//	opts := ClientOptions{
//	    ApiKey:    "your-api-key",
//	    ApiSecret: "your-secret",
//	    Timeout:   5 * time.Second,
//	}
//
//	client, err := NewClient(opts)
//	if err != nil {
//	    panic(err)
//	}
//
//	// The client is now ready to call API methods:
//	info, _ := client.GetMarketInformation()
func NewClient(opts ClientOptions) (*Client, error) {
	client := &Client{
		BaseUrl: BaseUrl,
	}

	if opts.BaseUrl != "" {
		client.BaseUrl = opts.BaseUrl
	}

	if opts.ApiKey != "" {
		client.ApiKey = opts.ApiKey
	}

	if opts.ApiSecret != "" {
		client.ApiSecret = opts.ApiSecret
	}

	if opts.HttpClient != nil {
		client.HttpClient = opts.HttpClient
	} else {
		client.HttpClient = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	return client, nil
}

// assertAuth validates that the client is currently authenticated by checking
// whether an ApiKey is available.
//
// Parameters:
//   - client: A pointer to the Client instance.
//
// Returns:
//   - nil if the client contains a non-empty ApiKey.
//   - A *GoTabdealError if ApiKey is empty.
//
// Behavior:
//   - This function does not perform any network I/O.
//   - It is used internally before making authenticated requests.
//
// Errors:
//   - "api key is empty" when ApiKey is missing.
//   - "api secret is empty" when ApiSecret is missing.
//
// Example:
//
//	if err := assertAuth(client); err != nil {
//	    return err
//	}
func assertAuth(client *Client) error {
	if client.ApiKey == "" {
		return &GoTabdealError{
			Message: "api key is empty",
			Err:     nil,
		}
	}

	if client.ApiSecret == "" {
		return &GoTabdealError{
			Message: "api secret is empty",
			Err:     nil,
		}
	}

	return nil
}

// createApiURI constructs a fully qualified Tabdeal API URL using the client's
// base URL, an optional version prefix, and the raw endpoint path.
//
// Parameters:
//   - endpoint: The API endpoint (MUST begin with a leading slash), for example:
//     "/market/orders/add"
//     "/options"
//     "/orderbook/BTCUSDT/"
//   - version: Optional API version string such as "v2" or "v3".
//     If empty, no version segment is inserted.
//
// Returns:
//   - A fully qualified URL as:
//   - Without version:  {BaseUrl}{endpoint}
//   - With version:     {BaseUrl}/{version}{endpoint}
//
// Behavior:
//   - BaseUrl must NOT have a trailing slash (e.g. "https://apiv2.Tabdeal.ir").
//   - Endpoint MUST begin with "/", and is appended as-is.
//   - Version MUST NOT begin with "/", the function prepends one automatically.
//
// Examples:
//
//	c.BaseUrl = "https://apiv2.Tabdeal.ir"
//
//	createApiURI("/market/stats", "")
//	→ "https://apiv2.Tabdeal.ir/market/stats"
//
//	createApiURI("/options", "v2")
//	→ "https://apiv2.Tabdeal.ir/v2/options"
//
//	createApiURI("/orderbook/BTCUSDT", "v3")
//	→ "https://apiv2.Tabdeal.ir/v3/orderbook/BTCUSDT"
func (c *Client) createApiURI(method string, endpoint string) string {
	if method == "GET" {
		return fmt.Sprintf("%s/r/api/%s%s", c.BaseUrl, c.Version, endpoint)
	}

	return fmt.Sprintf("%s/api/%s%s", c.BaseUrl, c.Version, endpoint)
}

// Request sends an HTTP request to a full Tabdeal URL and handles:
//   - optional authentication
//   - automatic TOTP header placement
//   - request serialization (JSON or query params)
//   - response deserialization
//   - structured API error handling
//
// Parameters:
//   - method: "GET" or "POST".
//   - url: Fully constructed URL (already includes version when needed).
//   - auth: Whether the request requires an Authorization header.
//   - otpRequired: Whether X-TOTP should be sent for this specific request.
//   - body: For GET, serialized into URL params; for POST, JSON-encoded.
//   - result: Optional pointer to the output struct into which JSON is unmarshaled.
//
// Returns:
//   - nil on success.
//   - *RequestError on network/encoding issues.
//   - *APIError when Tabdeal returns status != 2xx.
//
// Behavior:
//   - GET: struct → ?a=b&c=d via StructToURLParams.
//   - POST: struct → JSON in request body.
//   - If auth=true:
//   - handleAutoRefresh() is executed when AutoRefresh is enabled.
//   - assertAuth() ensures ApiKey is set.
//   - User-Agent and Authorization headers are required.
//   - X-TOTP header is added when otpRequired=true.
//   - On error HTTP status, parseErrorResponse() maps Tabdeal JSON error objects
//     into APIError (fields: status, code, message, detail).
//
// Dependencies:
//   - StructToURLParams
//   - assertAuth()
//   - handleAutoRefresh()
//   - parseErrorResponse()
//
// Errors:
//   - "failed to marshal request body"
//   - "failed to convert struct to URL params"
//   - "failed to send request"
//   - "failed to unmarshal response"
//   - APIError (status, code, message, detail)
//
// Example:
//
//	var res t.OrderStatus
//	err := client.Request("POST", url, true, true, params, &res)
//	if err != nil {
//	    return err
//	}
func (c *Client) Request(method string, url string, auth bool, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		var urlParams string
		if auth {
			bodyUpdate := u.WrapWithSignature(body, c.ApiSecret, time.Now().Unix()*1000)
			urlParams, err = u.StructToURLParams(bodyUpdate)
		} else {
			urlParams, err = u.StructToURLParams(body)
		}

		if err != nil {
			return &RequestError{
				GoTabdealError: GoTabdealError{
					Message: "failed to convert struct to URL params",
					Err:     err,
				},
				Operation: "preparing request parameters",
			}
		}
		url += "?" + urlParams
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return &RequestError{
			GoTabdealError: GoTabdealError{
				Message: "failed to create request",
				Err:     err,
			},
			Operation: "creating request",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	if auth {
		if err := assertAuth(c); err != nil {
			return &GoTabdealError{
				Message: "authentication validation failed",
				Err:     err,
			}
		}

		req.Header.Set("X-MBX-APIKEY", c.ApiKey)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return &RequestError{
			GoTabdealError: GoTabdealError{
				Message: "failed to send request",
				Err:     err,
			},
			Operation: "sending request",
		}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestError{
			GoTabdealError: GoTabdealError{
				Message: "failed to read response body",
				Err:     err,
			},
			Operation: "reading response",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseErrorResponse(resp.StatusCode, respBody)
	}

	if result != nil {
		if err = json.Unmarshal(respBody, result); err != nil {
			return &RequestError{
				GoTabdealError: GoTabdealError{
					Message: "failed to unmarshal response",
					Err:     err,
				},
				Operation: "parsing response",
			}
		}
	}

	return nil
}

// ApiRequest is a convenience wrapper that builds a Tabdeal API URL using
// createApiURI() and delegates the actual HTTP call to Request().
//
// Parameters:
//   - method: HTTP method ("GET", "POST").
//   - endpoint: The endpoint path, such as "/market/orders/add".
//   - version: Tabdeal version string ("v2", "v3"). May be empty.
//   - auth: Whether this call requires Authorization: Token <key>.
//   - otpRequired: Whether this endpoint requires X-TOTP.
//   - body: Struct for GET params or POST JSON body.
//   - result: Destination struct for response JSON.
//
// Returns:
//   - nil on success.
//   - See Request() for structured errors.
//
// Behavior:
//   - Constructs URL = BaseUrl + "/api/{version}/{endpoint}".
//   - Passes all fields to Request() unchanged.
//
// Example:
//
//	var stats t.Tickers
//	err := client.ApiRequest("GET", "/market/stats", "", false, false, params, &stats)
func (c *Client) ApiRequest(method, endpoint string, auth bool, body interface{}, result interface{}) error {
	url := c.createApiURI(method, endpoint)
	return c.Request(method, url, auth, body, result)
}

func (c *Client) ping() (bool, error) {
	err := c.ApiRequest("GET", "/ping", false, nil, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) GetServerTime() (*t.ServerTime, error) {
	var serverTime *t.ServerTime
	err := c.ApiRequest("GET", "/time", false, nil, &serverTime)
	if err != nil {
		return nil, err
	}
	return serverTime, nil
}

// GetMarketInformation retrieves full market metadata from Tabdeal,
// including symbol configurations, precision rules, filters, and
// trading permissions.
//
// Endpoint:
//
//	GET /r/api/v1/exchangeInfo
//
// Returns:
//   - []*t.MarketInformation: list of market definitions provided by Tabdeal.
//   - error on network or API failure.
//
// Behavior:
//   - No authentication required.
//   - Returns exchange-wide configuration for all markets.
//
// Example:
//
//	info, err := client.GetMarketInformation()
//	if err != nil { panic(err) }
//	fmt.Println(info[0].Symbol)
func (c *Client) GetMarketInformation() (*[]*t.MarketInformation, error) {
	var marketInfo *[]*t.MarketInformation
	err := c.ApiRequest("GET", "/exchangeInfo", false, nil, &marketInfo)
	if err != nil {
		return nil, err
	}
	return marketInfo, nil
}

// GetOrderBook returns the current depth (order book) for a specific market.
//
// Endpoint:
//
//	GET /r/api/v1/depth?symbol=SYMBOL
//
// Params:
//   - params.Symbol: market symbol (BTCUSDT, BTCIRT, etc.)
//
// Returns:
//   - *t.OrderBook including raw bids and asks.
//   - error on failure.
//
// Behavior:
//   - No authentication required.
//   - Depth is returned in aggregated form: [][]string.
//
// Example:
//
//	book, _ := client.GetOrderBook(t.GetOrderBookParams{Symbol: "BTCUSDT"})
//	fmt.Println(book.Bids[0])
func (c *Client) GetOrderBook(params t.GetOrderBookParams) (*t.OrderBook, error) {
	var orderBook *t.OrderBook
	err := c.ApiRequest("GET", "/depth", false, params, &orderBook)
	if err != nil {
		return nil, err
	}
	return orderBook, nil
}

// GetRecentTrades retrieves most recent trades for a market.
//
// Endpoint:
//
//	GET /r/api/v1/trades?symbol=SYMBOL&limit=...
//
// Params:
//   - Symbol (required)
//   - Limit (optional)
//
// Returns:
//   - []*t.Trade sorted newest → oldest.
//   - error on API or network failure.
//
// Behavior:
//   - No authentication required.
//
// Example:
//
//	trades, _ := client.GetRecentTrades(t.GetRecentTradesParams{Symbol: "BTCUSDT"})
//	fmt.Println(trades[0].Price)
func (c *Client) GetRecentTrades(params t.GetRecentTradesParams) (*[]*t.Trade, error) {
	var trades *[]*t.Trade
	err := c.ApiRequest("GET", "/trades", false, params, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}

// GetWallets retrieves funding wallet balances for the authenticated user.
//
// Endpoint:
//
//	GET /api/v1/get-funding-asset?asset=...
//
// Params:
//   - Asset (optional): filter by a specific asset.
//
// Authentication:
//   - Required. Uses X-MBX-APIKEY header.
//
// Returns:
//   - []*t.Wallet containing asset, free, and freeze amounts.
//   - error if authentication is missing or API responds with error.
//
// Example:
//
//	balances, _ := client.GetWallets(t.GetWalletParams{Asset: "USDT"})
//	fmt.Println(balances[0].Free)
func (c *Client) GetWallets(params t.GetWalletParams) (*[]*t.Wallet, error) {
	var wallets *[]*t.Wallet
	err := c.ApiRequest("GET", "/get-funding-asset", true, params, &wallets)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// CreateOrder submits a new spot order to Tabdeal.
//
// Endpoint:
//
//	POST /api/v1/order
//
// Params (t.CreateOrderParams):
//   - Side ("BUY"/"SELL")
//   - Type ("LIMIT", "MARKET")
//   - Quantity (required)
//   - Price (for limit orders)
//   - StopPrice (for stop orders)
//   - Symbol/TabdealSymbol (from BaseSymbolParams)
//
// Authentication:
//   - Required. Signs parameters using API secret.
//
// Returns:
//   - *t.CreateOrderResponse with full order details and fills.
//   - error on failure.
//
// Example:
//
//	resp, _ := client.CreateOrder(t.CreateOrderParams{
//	    Symbol: "BTCUSDT",
//	    Side: "BUY",
//	    Type: "LIMIT",
//	    Quantity: 0.01,
//	    Price: 950000000,
//	})
func (c *Client) CreateOrder(params t.CreateOrderParams) (*t.CreateOrderResponse, error) {
	var createOrderResponse *t.CreateOrderResponse
	err := c.ApiRequest("POST", "/order", true, params, &createOrderResponse)
	if err != nil {
		return nil, err
	}
	return createOrderResponse, nil
}

// CancelOrder cancels a single active order.
//
// Endpoint:
//
//	DELETE /api/v1/order?symbol=...&orderId=...
//
// Params:
//   - Symbol
//   - OrderId OR OrigClientOrderId
//
// Authentication:
//   - Required. Signed request.
//
// Returns:
//   - *t.CancelOrderResponse containing BaseOrderResponse
//   - error on API or network failure.
//
// Example:
//
//	client.CancelOrder(t.CancelOrderParams{
//	    Symbol: "BTCUSDT",
//	    OrderId: 1234567,
//	})
func (c *Client) CancelOrder(params t.CancelOrderParams) (*t.CancelOrderResponse, error) {
	var cancelOrderStatus *t.CancelOrderResponse
	err := c.ApiRequest("DELETE", "/order", true, params, &cancelOrderStatus)
	if err != nil {
		return nil, err
	}
	return cancelOrderStatus, nil
}

// CancelOrderBulk cancels all open orders for a symbol.
//
// Endpoint:
//
//	DELETE /api/v1/openOrders?symbol=...
//
// Params:
//   - Symbol (from BaseSymbolParams)
//
// Authentication:
//   - Required.
//
// Returns:
//   - []*t.CancelOrderResponse for all canceled orders.
//   - error on API or network failure.
//
// Example:
//
//	resp, _ := client.CancelOrderBulk(t.CancelOrderBulkParams{
//	    Symbol: "BTCUSDT",
//	})
func (c *Client) CancelOrderBulk(params t.CancelOrderBulkParams) (*[]*t.CancelOrderResponse, error) {
	var cancelOrderBulkStatus *[]*t.CancelOrderResponse
	err := c.ApiRequest("DELETE", "/openOrders", true, params, &cancelOrderBulkStatus)
	if err != nil {
		return nil, err
	}
	return cancelOrderBulkStatus, nil
}

// GetOrdersHistory retrieves historical orders for the authenticated user.
//
// Endpoint:
//
//	GET /api/v1/allOrders?symbol=...&startTime=...&endTime=...&limit=...
//
// Params:
//   - Symbol
//   - StartTime / EndTime (optional)
//   - Limit (optional)
//
// Authentication:
//   - Required.
//
// Returns:
//   - []*t.BaseOrderResponse
//   - error
//
// Example:
//
//	history, _ := client.GetOrdersHistory(t.GetUserOrdersHistoryParams{
//	    Symbol: "BTCUSDT",
//	    Limit: 50,
//	})
func (c *Client) GetOrdersHistory(params t.GetUserOrdersHistoryParams) (*[]*t.BaseOrderResponse, error) {
	var orders *[]*t.BaseOrderResponse
	err := c.ApiRequest("GET", "/allOrders", true, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOpenOrders retrieves all currently-open orders for a market.
//
// Endpoint:
//
//	GET /api/v1/openOrders?symbol=...
//
// Params:
//   - Symbol (from BaseSymbolParams)
//
// Authentication:
//   - Required.
//
// Returns:
//   - []*t.BaseOrderResponse
//
// Example:
//
//	open, _ := client.GetOpenOrders(t.GetOpenOrdersParams{Symbol: "BTCUSDT"})
func (c *Client) GetOpenOrders(params t.GetOpenOrdersParams) (*[]*t.BaseOrderResponse, error) {
	var orders *[]*t.BaseOrderResponse
	err := c.ApiRequest("GET", "/openOrders", true, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderStatus retrieves the full status of a specific order.
//
// Endpoint:
//
//	GET /api/v1/order?symbol=...&orderId=... OR &origClientOrderId=...
//
// Params:
//   - OrderId or OrigClientOrderId
//
// Authentication:
//   - Required.
//
// Returns:
//   - *t.OrderStatusResponse
//
// Example:
//
//	st, _ := client.GetOrderStatus(t.GetOrderStatusParams{OrderId: 1234})
func (c *Client) GetOrderStatus(params t.GetOrderStatusParams) (*t.OrderStatusResponse, error) {
	var orders *t.OrderStatusResponse
	err := c.ApiRequest("GET", "/order", true, params, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetUserTrades returns the authenticated user's trade history for a symbol.
//
// Endpoint:
//
//	GET /api/v1/myTrades?symbol=...&fromId=...
//
// Params:
//   - Symbol
//   - FromId (pagination, optional)
//
// Authentication:
//   - Required.
//
// Returns:
//   - []*t.UserTradeResponse
//   - error
//
// Example:
//
//	trades, _ := client.GetUserTrades(t.GetUserTradesParams{
//	    Symbol: "BTCUSDT",
//	})
func (c *Client) GetUserTrades(params t.GetUserTradesParams) (*[]*t.UserTradeResponse, error) {
	var trades *[]*t.UserTradeResponse
	err := c.ApiRequest("GET", "/myTrades", true, params, &trades)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
