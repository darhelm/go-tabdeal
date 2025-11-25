# Tabdeal Go SDK — Full Examples

This document includes practical examples for every major SDK function.

---

# Initialize Client

```go
client, err := tabdeal.NewClient(tabdeal.ClientOptions{
    ApiKey:    "API_KEY",
    ApiSecret: "SECRET",
    Timeout:   5 * time.Second,
})
```

---

# Market Information

## Get Market Information

```go
info, err := client.GetMarketInformation()
fmt.Println((*info)[0].Symbol)
```

## Get Order Book

```go
ob, err := client.GetOrderBook(types.GetOrderBookParams{
    Symbol: "BTCIRT",
})
fmt.Println(ob.Asks[0], ob.Bids[0])
```

## Get Recent Trades

```go
recent, err := client.GetRecentTrades(types.GetRecentTradesParams{
    Symbol: "BTCIRT",
})
fmt.Println((*recent)[0])
```

---

# Wallet Operations

## Get Wallets

```go
wallets, err := client.GetWallets(types.GetWalletParams{
    Asset: "USDT",
})
fmt.Println((*wallets)[0].Free)
```

---

# Trading

## Create Order

```go
createResp, err := client.CreateOrder(types.CreateOrderParams{
    Execution:   "limit",
    Type:        "buy",
    SrcCurrency: "btc",
    DstCurrency: "usdt",
    Amount:      "0.01",
    Price:       "1500000000",
})
fmt.Println(createResp.Status)
```

## Cancel Order

```go
cancelResp, err := client.CancelOrder(types.CancelOrderParams{
    Id: 999,
})
```

## Bulk Cancel

```go
bulkResp, err := client.CancelOrderBulk(types.CancelOrderBulkParams{
    Hours: 12,
})
```

## Orders History

```go
hist, err := client.GetOrdersHistory(types.GetUserOrdersHistoryParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
fmt.Println((*hist)[0])
```

## Open Orders

```go
open, err := client.GetOpenOrders(types.GetOpenOrdersParams{})
fmt.Println((*open)[0])
```

## Order Status

```go
st, err := client.GetOrderStatus(types.GetOrderStatusParams{
    Id: 777,
})
fmt.Println(st)
```

---

# User Trades

```go
trades, err := client.GetUserTrades(types.GetUserTradesParams{
    SrcCurrency: "btc",
    DstCurrency: "usdt",
})
fmt.Println((*trades)[0])
```

---

# Error Handling

```go
_, err := client.GetMarketInformation()
if err != nil {
    if apiErr, ok := err.(*tabdeal.APIError); ok {
        fmt.Println(apiErr.Code)
        fmt.Println(apiErr.Message)
        fmt.Println(apiErr.Detail)
    }
}
```