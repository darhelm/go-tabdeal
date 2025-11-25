# Go Tabdeal

[![Go Reference](https://pkg.go.dev/badge/github.com/darhelm/go-tabdeal.svg)](https://pkg.go.dev/github.com/darhelm/go-tabdeal)
[![Go Report Card](https://goreportcard.com/badge/github.com/darhelm/go-tabdeal)](https://goreportcard.com/report/github.com/darhelm/go-tabdeal)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/darhelm/go-tabdeal)](https://golang.org/dl/)

A clean, strongly typed, and fully documented Go SDK for interacting with the **Tabdeal** exchange API.

This SDK provides strict typed models, request signing, public/private endpoints, and predictable behavior with minimal abstraction.

## Disclaimer
This SDK is **unofficial**. Use at your own risk.

## Features
- Full support for Tabdeal public & authenticated endpoints
- Strongly typed request/response structs
- Simple API key + secret authentication
- Request signing (HMAC-SHA256)
- Order placement, cancellation, bulk cancellation
- Wallets, trades, order history
- Order book & recent trades
- Fully structured error handling (`APIError`, `RequestError`)

## Installation
```bash
go get github.com/darhelm/go-tabdeal
```

## Quick Start
```go
package main

import (
    "fmt"
    "time"
    tabdeal "github.com/darhelm/go-tabdeal"
)

func main() {
    client, err := tabdeal.NewClient(tabdeal.ClientOptions{
        ApiKey:    "YOUR_KEY",
        ApiSecret: "YOUR_SECRET",
        Timeout:   5 * time.Second,
    })
    if err != nil {
        panic(err)
    }

    info, err := client.GetMarketInformation()
    if err != nil {
        panic(err)
    }

    fmt.Println((*info)[0].Symbol)
}
```

## Documentation
- SDK Reference: https://pkg.go.dev/github.com/darhelm/go-tabdeal
- Tabdeal API Docs: https://docs.tabdeal.org/
- Full examples: `EXAMPLES.md`

---

# Examples

### Create Client
```go
client, err := tabdeal.NewClient(tabdeal.ClientOptions{
    ApiKey:    "KEY",
    ApiSecret: "SECRET",
    Timeout:   5 * time.Second,
})
```

### Market Information
```go
info, err := client.GetMarketInformation()
fmt.Println((*info)[0].Symbol)
```

### Order Book
```go
ob, err := client.GetOrderBook(types.GetOrderBookParams{
    Symbol: "BTCIRT",
})
fmt.Println(ob.Asks[0], ob.Bids[0])
```

### Recent Trades
```go
trades, err := client.GetRecentTrades(types.GetRecentTradesParams{
    Symbol: "BTCIRT",
})
```

### Wallets
```go
wallets, err := client.GetWallets(types.GetWalletParams{
    Asset: "USDT",
})
```

### Create Order
```go
resp, err := client.CreateOrder(types.CreateOrderParams{
    Symbol:   "BTCIRT",
    Side:     "BUY",
    Type:     "LIMIT",
    Quantity: 0.01,
    Price:    1500000000,
})
```

### Cancel Order
```go
cancel, err := client.CancelOrder(types.CancelOrderParams{
    Symbol:  "BTCIRT",
    OrderId: 1234,
})
```

### Bulk Cancel
```go
bulk, err := client.CancelOrderBulk(types.CancelOrderBulkParams{
    Symbol: "BTCIRT",
})
```

### Order History
```go
history, err := client.GetOrdersHistory(types.GetUserOrdersHistoryParams{
    Symbol: "BTCIRT",
})
```

### Open Orders
```go
open, err := client.GetOpenOrders(types.GetOpenOrdersParams{
    Symbol: "BTCIRT",
})
```

### Order Status
```go
st, err := client.GetOrderStatus(types.GetOrderStatusParams{
    OrderId: 1234,
})
```

### User Trades
```go
trades, err := client.GetUserTrades(types.GetUserTradesParams{
    Symbol: "BTCIRT",
})
```

### Error Handling
```go
if err != nil {
    if apiErr, ok := err.(*tabdeal.APIError); ok {
        fmt.Println(apiErr.Code, apiErr.Message, apiErr.Detail)
    }
}
```

## License
MIT License.