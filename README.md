# go-hyperliquid
A Golang SDK for the Hyperliquid API.

# API reference
- [Hyperliquid](https://app.hyperliquid.xyz/)
- [Hyperliquid API docs](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api)
- [Hyperliquid official Python SDK](https://github.com/hyperliquid-dex/hyperliquid-python-sdk)

# How to install?
```
go get github.com/chainswatch/go-hyperliquid
```

# Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/chainswatch/go-hyperliquid.svg)](https://pkg.go.dev/github.com/chainswatch/go-hyperliquid)


# Quick start
```
package main

import (
	"log"

        "github.com/chainswatch/go-hyperliquid"
)

func main() {
	hyperliquidClient := hyperliquid.NewHyperliquid(&hyperliquid.HyperliquidClientConfig{
		IsMainnet:      true,
		AccountAddress: "0x12345",   // Main address of the Hyperliquid account that you want to use
		PrivateKey:     "abc1234",   // Private key of the account or API private key from Hyperliquid
	})

	// Get balances
	res, err := hyperliquidClient.GetAccountState()
	if err != nil {
		log.Print(err)
	}
	log.Printf("GetAccountState(): %+v", res)
}
```

# Running tests

Integration tests require access to a funded Hyperliquid account. Provide the credentials via environment variables `TEST_ADDRESS` and `TEST_PRIVATE_KEY`. For convenience you can copy `.test.env.example` to `.test.env` at the repository root and populate these variables:

```
TEST_ADDRESS=0xabc123...
TEST_PRIVATE_KEY=...
```

The test suite automatically loads this file if present. Avoid committing any real credentials to source control.
Environment files like `.env` and `.test.env` are ignored by git.

If these variables are missing when the tests run the suite will fail immediately.

Run the tests with:

```
make test
```

The Makefile uses `go test -count=1` so results are not cached between runs.


