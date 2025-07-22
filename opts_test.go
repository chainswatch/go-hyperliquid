package hyperliquid

import (
	"net/http"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestClientOptions_WithHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      false,
		PrivateKey:     "",
		AccountAddress: "",
	}, WithHTTPClient(customClient))

	if hl.ExchangeAPI.httpClient != customClient {
		t.Errorf("ExchangeAPI.httpClient = %p, want %p", hl.ExchangeAPI.httpClient, customClient)
	}

	if hl.InfoAPI.httpClient != customClient {
		t.Errorf("InfoAPI.httpClient = %p, want %p", hl.InfoAPI.httpClient, customClient)
	}

	if hl.ExchangeAPI.httpClient.Timeout != 30*time.Second {
		t.Errorf("ExchangeAPI.httpClient.Timeout = %v, want %v", hl.ExchangeAPI.httpClient.Timeout, 30*time.Second)
	}
}

func TestClientOptions_WithLogger(t *testing.T) {
	customLogger := log.New()
	customLogger.SetLevel(log.ErrorLevel)

	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      false,
		PrivateKey:     "",
		AccountAddress: "",
	}, WithLogger(customLogger))

	if hl.ExchangeAPI.Logger != customLogger {
		t.Errorf("ExchangeAPI.Logger = %p, want %p", hl.ExchangeAPI.Logger, customLogger)
	}

	if hl.InfoAPI.Logger != customLogger {
		t.Errorf("InfoAPI.Logger = %p, want %p", hl.InfoAPI.Logger, customLogger)
	}

	if hl.ExchangeAPI.Logger.Level != log.ErrorLevel {
		t.Errorf("ExchangeAPI.Logger.Level = %v, want %v", hl.ExchangeAPI.Logger.Level, log.ErrorLevel)
	}
}

func TestClientOptions_MultipleOptions(t *testing.T) {
	customClient := &http.Client{
		Timeout: 45 * time.Second,
	}
	customLogger := log.New()
	customLogger.SetLevel(log.WarnLevel)

	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      true,
		PrivateKey:     "",
		AccountAddress: "",
	}, WithHTTPClient(customClient), WithLogger(customLogger))

	if hl.ExchangeAPI.httpClient != customClient {
		t.Errorf("ExchangeAPI.httpClient = %p, want %p", hl.ExchangeAPI.httpClient, customClient)
	}

	if hl.ExchangeAPI.Logger != customLogger {
		t.Errorf("ExchangeAPI.Logger = %p, want %p", hl.ExchangeAPI.Logger, customLogger)
	}

	if hl.InfoAPI.httpClient != customClient {
		t.Errorf("InfoAPI.httpClient = %p, want %p", hl.InfoAPI.httpClient, customClient)
	}

	if hl.InfoAPI.Logger != customLogger {
		t.Errorf("InfoAPI.Logger = %p, want %p", hl.InfoAPI.Logger, customLogger)
	}
}

func TestClientOptions_NoOptions(t *testing.T) {
	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      false,
		PrivateKey:     "",
		AccountAddress: "",
	})

	if hl.ExchangeAPI.httpClient != http.DefaultClient {
		t.Errorf("ExchangeAPI.httpClient = %p, want %p", hl.ExchangeAPI.httpClient, http.DefaultClient)
	}

	if hl.InfoAPI.httpClient != http.DefaultClient {
		t.Errorf("InfoAPI.httpClient = %p, want %p", hl.InfoAPI.httpClient, http.DefaultClient)
	}

	if hl.ExchangeAPI.Logger == nil {
		t.Error("ExchangeAPI.Logger should not be nil")
	}

	if hl.InfoAPI.Logger == nil {
		t.Error("InfoAPI.Logger should not be nil")
	}
}

func TestClientOptions_BackwardCompatibility(t *testing.T) {
	// Test that existing code without options still works
	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      true,
		PrivateKey:     "",
		AccountAddress: "0x1234567890abcdef1234567890abcdef12345678",
	})

	if hl.AccountAddress() != "0x1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("AccountAddress = %s, want %s", hl.AccountAddress(), "0x1234567890abcdef1234567890abcdef12345678")
	}

	if !hl.IsMainnet() {
		t.Error("IsMainnet() should return true")
	}

	// Test with nil config
	hlNil := NewHyperliquid(nil)
	if !hlNil.IsMainnet() {
		t.Error("Default config should use mainnet")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	customClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	customLogger := log.New()
	customLogger.SetLevel(log.InfoLevel)

	client := NewClient(false, WithHTTPClient(customClient), WithLogger(customLogger))

	if client.httpClient != customClient {
		t.Errorf("Client.httpClient = %p, want %p", client.httpClient, customClient)
	}

	if client.Logger != customLogger {
		t.Errorf("Client.Logger = %p, want %p", client.Logger, customLogger)
	}

	if client.IsMainnet() {
		t.Error("Client.IsMainnet() should return false")
	}
}

func TestNewExchangeAPI_WithOptions(t *testing.T) {
	customClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	api := NewExchangeAPI(true, WithHTTPClient(customClient))

	if api.httpClient != customClient {
		t.Errorf("ExchangeAPI.httpClient = %p, want %p", api.httpClient, customClient)
	}

	if !api.IsMainnet() {
		t.Error("ExchangeAPI.IsMainnet() should return true")
	}
}

func TestNewInfoAPI_WithOptions(t *testing.T) {
	customLogger := log.New()
	customLogger.SetLevel(log.FatalLevel)

	api := NewInfoAPI(false, WithLogger(customLogger))

	if api.Logger != customLogger {
		t.Errorf("InfoAPI.Logger = %p, want %p", api.Logger, customLogger)
	}

	if api.IsMainnet() {
		t.Error("InfoAPI.IsMainnet() should return false")
	}
}

func TestClientOptions_WithDebug(t *testing.T) {
	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      false,
		PrivateKey:     "",
		AccountAddress: "",
	}, WithDebug(true))

	if !hl.ExchangeAPI.Debug {
		t.Error("ExchangeAPI.Debug should be true")
	}

	if !hl.InfoAPI.Debug {
		t.Error("InfoAPI.Debug should be true")
	}

	// Test with debug false
	hlDebugOff := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      false,
		PrivateKey:     "",
		AccountAddress: "",
	}, WithDebug(false))

	if hlDebugOff.ExchangeAPI.Debug {
		t.Error("ExchangeAPI.Debug should be false")
	}

	if hlDebugOff.InfoAPI.Debug {
		t.Error("InfoAPI.Debug should be false")
	}
}

func TestClientOptions_WithDebugAndOtherOptions(t *testing.T) {
	customClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	hl := NewHyperliquid(&HyperliquidClientConfig{
		IsMainnet:      true,
		PrivateKey:     "",
		AccountAddress: "",
	}, WithHTTPClient(customClient), WithDebug(true))

	if hl.ExchangeAPI.httpClient != customClient {
		t.Errorf("ExchangeAPI.httpClient = %p, want %p", hl.ExchangeAPI.httpClient, customClient)
	}

	if !hl.ExchangeAPI.Debug {
		t.Error("ExchangeAPI.Debug should be true")
	}

	if !hl.InfoAPI.Debug {
		t.Error("InfoAPI.Debug should be true")
	}
}

