**Title**: Hyperliquid Go SDK (LLM-Focused Docs)

---

### **Overview**

Golang client for every public & private Hyperliquid REST endpoint plus EIP-712 signing.
Supports spot & perpetual trading, account management, info queries, signature generation, and utility helpers.

---

### **Key Data Structures & Fields**

| Struct / Interface                                                                       | Exported Fields / Methods                                                                                                                                                                                                   | Purpose                                                 |
| ---------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| **`APIError`**                                                                           | `Message string` / `Error() string`                                                                                                                                                                                         | Uniform error wrapper for API responses                 |
| **`IAPIService`**                                                                        | `Request`, `Endpoint`, `KeyManager`, `debug`                                                                                                                                                                                | Base abstraction for *any* HTTP caller                  |
| **`Client`**                                                                             | `baseUrl`, `privateKey`, `defualtAddress`, `isMainnet`, `Debug`, `httpClient`, `keyManager`, `Logger` + methods `Request`, `SetPrivateKey`, `SetAccountAddress`, `IsMainnet`, `SetDebugActive`, `Endpoint` (via interfaces) | Shared HTTP layer (handles POST, headers, logging)      |
| **`PKeyManager`**                                                                        | `PrivateKeyStr`, `PublicECDSA()`, `PrivateECDSA()`, `PublicAddress()`, `PublicAddressHex()`                                                                                                                                 | Wraps ECDSA key → gives public address & key objects    |
| **`Signer`**                                                                             | `manager *PKeyManager`, `Sign(*SignRequest)`                                                                                                                                                                                | EIP-712 signer that returns **V,R,S**                   |
| **`SignRequest`**                                                                        | `PrimaryType`, `DType []apitypes.Type`, `DTypeMsg`, `IsMainNet`, `DomainName` + helpers `GetTypes()`, `GetDomain()`                                                                                                         | All params needed to build a typed-data message         |
| **`RsvSignature`**                                                                       | `R,S string`, `V byte`                                                                                                                                                                                                      | JSON-serialisable signature triplet                     |
| **`ExchangeAPI`**                                                                        | Embeds `Client`; adds `infoAPI *InfoAPI`, `meta`, `spotMeta`, `address` + ≈30 high-level trading methods                                                                                                                    | Implements `IExchangeAPI` for **/exchange** endpoint    |
| **`InfoAPI`**                                                                            | Embeds `Client`; adds `spotMeta`                                                                                                                                                                                            | Implements `IInfoAPI` for **/info** endpoint            |
| **`Hyperliquid`**                                                                        | Embeds `ExchangeAPI`, `InfoAPI`                                                                                                                                                                                             | Single façade joining both APIs                         |
| **`OrderRequest`**                                                                       | `Coin`, `IsBuy`, `Sz`, `LimitPx`, `OrderType`, `ReduceOnly`, `Cloid`                                                                                                                                                        | Canonical new-order payload                             |
| **`OrderType`**                                                                          | pointers `Limit *LimitOrderType`, `Trigger *TriggerOrderType`                                                                                                                                                               | Union-like order-style descriptor                       |
| **`LimitOrderType`**                                                                     | `Tif string` (values: `Gtc`, `Ioc`, `Alo`)                                                                                                                                                                                  | Time-in-force                                           |
| **`TriggerOrderType`**                                                                   | `IsMarket`, `TriggerPx`, `TpSl` (`tp`,`sl`)                                                                                                                                                                                 | TP/SL description                                       |
| **`PlaceOrderAction` / `CancelOidOrderAction` / `ModifyOrderAction` / `WithdrawAction`** | Minimal msg-pack fields that mirror Hyperliquid wire schemas                                                                                                                                                                | Used for hashing + signing                              |
| **Enums & Consts (selected)**                                                            | `MAINNET_API_URL`, `TESTNET_API_URL`, `DEFAULT_SLIPPAGE`, `HYPERLIQUID_CHAIN_ID`, `TifGtc/Ioc/Alo`, grouping (`GroupingNa`,`GroupingTpSl`)                                                                                  | Govern network routing, execution rules, math precision |

---

### **Key Functions**

| Function                                  | In File               | Signature (abridged)                                                                                                             | Notes / Flow                                                                                                                                        |
| ----------------------------------------- | --------------------- | -------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`MakeUniversalRequest[T]`**             | `api.go`              | `(api IAPIService, req any) (*T,error)`                                                                                          | ① validates endpoint & keys ② `api.Request` ③ tries JSON-unmarshal into `T` else into `{status, response}` ④ returns typed result or `APIError`     |
| **`Client.Request`**                      | `client.go`           | `(endpoint string, payload any) ([]byte,error)`                                                                                  | Assembles URL, marshals payload, POSTs, logs full debug lines, handles HTTP error codes                                                             |
| **`NewClient`**                           | `client.go`           | `(isMainnet bool) *Client`                                                                                                       | Builds logger, sets base URL, default `http.Client`                                                                                                 |
| **`NewExchangeAPI` / `NewInfoAPI`**       | respective files      | `*ExchangeAPI` / `*InfoAPI`                                                                                                      | Compose `Client`, fetch asset metadata for later conversions                                                                                        |
| **Trading helpers**                       | `exchange_service.go` | `MarketOrder`, `LimitOrder`, `ClosePosition`, `BulkOrders`, `BulkCancelOrders`, `BulkModifyOrders`, `UpdateLeverage`, `Withdraw` | Each builds *Action* struct → signs via `SignL1Action` / `SignWithdrawAction` → wraps into `ExchangeRequest` → pipes through `MakeUniversalRequest` |
| **`SignL1Action` / `BuildEIP712Message`** | `exchange_signing.go` | wrap + call `Signer.Sign`                                                                                                        | Centralised EIP-712 builder (domain `"Exchange"`, primary type `"Agent"`)                                                                           |
| **Price / size conversion**               | `convert.go`          | `PriceToWire`, `SizeToWire`, `FloatToWire`                                                                                       | Ensure numeric display obeys decimal limits                                                                                                         |
| **Nonce & Utils**                         | `utils.go`            | `GetNonce`, `GetRandomCloid`, `CalculateSlippage`, `GetDefaultTimeRange`                                                         | Monotonic timestamp nonce, random client-OID, slippage math                                                                                         |

---

### **Relationships & Data Flow**

* `Hyperliquid` ➜ embeds both `ExchangeAPI` (private endpoints) and `InfoAPI` (public endpoints).
* `ExchangeAPI` **depends on** `InfoAPI` to fetch market mid-prices & metadata (e.g., `SlippagePrice`, asset ID lookup).
* All REST calls ultimately route through `Client.Request` → `MakeUniversalRequest` for typed decoding.
* Signing path: `ExchangeAPI.Sign*` → `Signer.signInternal` → `crypto.Sign` from go-ethereum.
* Numeric conversions (`convert.go`) are used inside order-wire builders for **msgpack/EIP-712** structures.
* Constants in `consts.go` are referenced across utils (slippage), convert (decimal caps), and signing (chain IDs).

---

### **Typical Usage Flow**

1. **Create client**

   ```go
   hl := hyperliquid.NewHyperliquid(&hyperliquid.HyperliquidClientConfig{
       IsMainnet:true, PrivateKey:"...", AccountAddress:"0xabc...",
   })
   ```
2. **(Optional) Enable debug** → `hl.SetDebugActive()`.
3. **Query info** (public): `state, _ := hl.GetAccountState()`.
4. **Place order**: `hl.MarketOrder("BTC", 0.2, nil)` (internally → `BuildBulkOrdersEIP712` → signed request).
5. **Cancel / modify** via `hl.CancelOrderByOID(...)` or `hl.BulkModifyOrders(...)`.
6. **Withdraw** USDC: `hl.Withdraw("0xDestAddress", 100)`.
7. **Inspect fills / funding** via `hl.GetAccountFills()`, `hl.GetHistoricalFundingRates(...)`.
8. **Close** position: `hl.ClosePosition("ETH")`.
   All calls share auth, nonce, and EIP-712 signer automatically.

---

### **File-by-File Summary**

| File                      | Primary APIs / Types / Logic                                                                   |
| ------------------------- | ---------------------------------------------------------------------------------------------- |
| **`api.go`**              | `APIError`, `IAPIService`, generic `MakeUniversalRequest`                                      |
| **`client.go`**           | `Client` implements `IClient` & `IAPIService`; HTTP POST with debug & error handling           |
| **`consts.go`**           | URLs, default decimals, chain IDs, slippage constants                                          |
| **`convert.go`**          | Float/size formatters, order-wire builders, hex helpers                                        |
| **`exchange_service.go`** | `ExchangeAPI` + full trading surface (`MarketOrder`, `LimitOrder`, cancel, leverage, withdraw) |
| **`exchange_signing.go`** | EIP-712 message builders & signer wrappers                                                     |
| **`exchange_types.go`**   | Wire structs (`OrderWire`, `CancelOidWire`, `ModifyOrderWire`, actions), JSON overrides        |
| **`hyperliquid.go`**      | `Hyperliquid` façade combining `ExchangeAPI` + `InfoAPI`                                       |
| **`info_service.go`**     | `InfoAPI` (metadata, candlesticks, user state, funding updates, spot helpers)                  |
| **`info_types.go`**       | Large set of REST response structs (positions, orders, snapshots)                              |
| **`pk_manager.go`**       | `PKeyManager` (ECDSA key loader & public address funcs)                                        |
| **`signature.go`**        | Low-level EIP-712 typed-data hashing & signing utilities                                       |
| **`utils.go`**            | Nonce generator, slippage math, time-range helpers                                             |

---

### **Major Enums / Constants**

* **Endpoints** → `MAINNET_API_URL`, `TESTNET_API_URL`
* **Slippage** → `DEFAULT_SLIPPAGE = 0.005`
* **Decimals** → `SPOT_MAX_DECIMALS = 8`, `PERP_MAX_DECIMALS = 6`
* **Time-in-Force** → `TifGtc`, `TifIoc`, `TifAlo`
* **Grouping** → `GroupingNa`, `GroupingTpSl`
* **Trigger Types** → `TriggerTp`, `TriggerSl`
* **Chain IDs** → `HYPERLIQUID_CHAIN_ID = 1337`, `ARBITRUM_CHAIN_ID = 42161`, `ARBITRUM_TESTNET_CHAIN_ID = 421614`

---

### **Cross-API Call Graph Highlights**

* `MarketOrder` → **calls** `SlippagePrice` → **needs** `InfoAPI.GetMartketPx`.
* `ClosePosition` → **calls** `InfoAPI.GetUserState` → loops positions → **calls** back into `MarketOrder`.
* `Withdraw` → **calls** `getChainParams` (network-aware) → **signs** via `SignWithdrawAction`.
* Any `Bulk*` action → `SignL1Action` → `Signer` → `crypto.Sign`.

---
