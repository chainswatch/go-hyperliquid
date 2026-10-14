# go-hyperliquid

**A Go SDK for the Hyperliquid API: market data, account state and order execution.**

[![Go Reference](https://pkg.go.dev/badge/github.com/chainswatch/go-hyperliquid.svg)](https://pkg.go.dev/github.com/chainswatch/go-hyperliquid)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-767BFF?style=flat)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.25-13161D?style=flat&logo=go)](go.mod)

Maintained by [Chainswatch](https://chainswatch.com). Forked from [Logarithm-Labs/go-hyperliquid](https://github.com/Logarithm-Labs/go-hyperliquid).

## Install

```sh
go get github.com/chainswatch/go-hyperliquid
```

## Quickstart

```go
package main

import (
	"log"
	"os"

	"github.com/chainswatch/go-hyperliquid"
)

func main() {
	client := hyperliquid.NewHyperliquid(&hyperliquid.HyperliquidClientConfig{
		IsMainnet:      true,
		AccountAddress: os.Getenv("HL_ADDRESS"),     // main address of the Hyperliquid account
		PrivateKey:     os.Getenv("HL_PRIVATE_KEY"), // account or API wallet private key; can be empty for public endpoints
	})

	// Mid prices for all coins (public endpoint)
	mids, err := client.GetAllMids()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BTC mid: %s", (*mids)["BTC"])

	// Account state (balances and positions)
	state, err := client.GetAccountState()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GetAccountState(): %+v", state)
}
```

Full API: [pkg.go.dev/github.com/chainswatch/go-hyperliquid](https://pkg.go.dev/github.com/chainswatch/go-hyperliquid).

## References

- [Hyperliquid](https://app.hyperliquid.xyz/)
- [Hyperliquid API docs](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api)
- [Hyperliquid official Python SDK](https://github.com/hyperliquid-dex/hyperliquid-python-sdk)

## Running tests

Integration tests require access to a funded Hyperliquid account. Provide the credentials via environment variables `TEST_ADDRESS` and `TEST_PRIVATE_KEY`. For convenience you can copy `.test.env.example` to `.test.env` at the repository root and populate these variables:

```
TEST_ADDRESS=0xabc123...
TEST_PRIVATE_KEY=...
```

The test suite automatically loads this file if present. Avoid committing any real credentials to source control.
Environment files like `.env` and `.test.env` are ignored by git.

If these variables are missing when the tests run the suite will fail immediately.

Run the unit tests with:

```
make test
```

Run the integration tests (requires `TEST_ADDRESS` and `TEST_PRIVATE_KEY`):

```
make integration
```

To run both sets together use:

```
make test-all
```

The Makefile uses `go test -count=1` so results are not cached between runs.

## License

Apache 2.0. See [LICENSE](LICENSE).

---

Built by [Chainswatch](https://chainswatch.com).
