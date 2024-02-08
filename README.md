# appstore

提供以下 API:
- app store server api
- [TODO] app store connect api

> 注意两者的 API key 是需要单独创建的

## Contents

- [Installation](#Installation)
- [Documentation](#Documentation)
- [API Examples](#API-Examples)
  - [app store server api](#app-store-server-api)
    - [Init Service](#init-service)
    - [LookupOrder](#LookupOrder)
    - [GetTransactionInfo](#GetTransactionInfo)
    - [GetTransactionHistory](#GetTransactionHistory)
    - [GetRefundHistory](#GetRefundHistory)

## Installation

- install
```bash
go get -u github.com/beanscc/appstore
```

- import it in you code
```go
import "github.com/beanscc/appstore"
```

## Documentation

[API documentation](https://pkg.go.dev/github.com/beanscc/appstore) for package

## API Examples

### app store server api

API 调用示例未给出的，请参见 `Test` [方法](appstoreserverapi/api_test.go)

#### Init Service

```go
config := Config{
    BundleID: `Your app’s bundle ID (Ex: “com.example.testbundleid2021”)`,
    Issuer:   `Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")`,
    KeyID:    `Your private key ID from App Store Connect (Ex: 2X9R4HXF34)`,
    PrivateKey: []byte(`-----BEGIN PRIVATE KEY-----
    YOUR PRIVATE KEY
    -----END PRIVATE KEY-----`),
    Timeout: 20 * time.Second,
}

token := appstoreserverapi.NewToken(&config)
service := appstoreserverapi.NewService(token)
```

#### LookupOrder

```go
token := appstoreserverapi.NewToken(&config)
service := appstoreserverapi.NewService(token).Debug(false)
customerOrderID := `MTV70QV5J9`
transactions, err := service.LookupOrder(context.Background(), customerOrderID)
if err != nil {
    log.Printf("[ERROR] service.LookupOrder failed. err:%v, customerOrderID:%s", err, customerOrderID)
}

for i, v := range transactions {
    // ....
}
```

#### GetTransactionInfo

```go
token := appstoreserverapi.NewToken(testConfig())
service := appstoreserverapi.NewService(token).Debug(false)
transactionID := `350001859400409`
got, err := service.GetTransactionInfo(context.Background(), transactionID)
if err != nil {
    log.Printf("[ERROR] Service.GetTransactionInfo failed. err:%v", err)
    return
}

log.Printf("Service.GetTransactionInfo: got:%#v", got)
```

#### GetTransactionHistory

```go
ctx := context.Background()
token := appstoreserverapi.NewToken(testConfig())
service := appstoreserverapi.NewService(token).Debug(false)
req := GetTransactionHistoryReq{
    TransactionID: `350001859400409`,
    Query:         nil,
}
got, err := service.GetTransactionHistory(ctx, &req)
if err != nil {
    log.Printf("[ERROR] Service.GetTransactionHistory failed. err:%v", err)
    return
}

total := 0
for n := 0; ; n++ {
    loop := n + 1
    log.Printf("[INFO] Service.GetTransactionHistory loop:%d base:%#v", loop, got.TransactionHistoryBase)
    for i, v := range got.Transactions {
        log.Printf("[INFO] Service.GetTransactionHistory loop:%d, got idx:%3d, v:%#v", loop, i, v)
        total++
    }

    if got.HasMore {
		// next page
        got, err = got.Next(context.Background())
    } else {
        break
    }
}

log.Printf("[INFO] total:%d", total)
```

#### GetRefundHistory

```go
ctx := context.Background()
token := appstoreserverapi.NewToken(testConfig())
service := appstoreserverapi.NewService(token).Debug(false)

transactionID := `140002007488219`
revision := ``
got, err := service.GetRefundHistory(ctx, transactionID, revision)
if err != nil {
    log.Printf("[ERROR] Service.GetRefundHistory failed. err:%v", err)
    return
}

total := 0
for n := 0; ; n++ {
    loop := n + 1
    log.Printf("[INFO] Service.GetRefundHistory loop:%d, got revision:%s", loop, got.Revision)
    for i, v := range got.SignedTransactions {
        transaction, err := v.GetTransaction()
        if err != nil {
            log.Printf("[ERROR] Service.GetRefundHistory v.GetTransaction failed. err:%v", err)
            return
        }
        log.Printf("[INFO] Service.GetRefundHistory loop:%d, got idx:%3d, v:%#v", loop, i, transaction)
        total++
    }

    if got.HasMore {
		// next page
        got, err = got.Next(context.Background())
    } else {
        break
    }
}
log.Printf("[INFO] total:%d", total)
```