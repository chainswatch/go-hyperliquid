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
