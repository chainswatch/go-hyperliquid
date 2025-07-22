package hyperliquid

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// ClientOption defines a function type for configuring the Hyperliquid client
type ClientOption func(*clientOptions)

// clientOptions holds all configurable options for the Hyperliquid client
type clientOptions struct {
	httpClient *http.Client
	logger     *log.Logger
	debug      bool
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(opts *clientOptions) {
		opts.httpClient = client
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *log.Logger) ClientOption {
	return func(opts *clientOptions) {
		opts.logger = logger
	}
}

// WithDebug sets the debug mode
func WithDebug(debug bool) ClientOption {
	return func(opts *clientOptions) {
		opts.debug = debug
	}
}

// getDefaultOptions returns the default client options
func getDefaultOptions() *clientOptions {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	return &clientOptions{
		httpClient: http.DefaultClient,
		logger:     logger,
		debug:      false,
	}
}

// applyOptions applies the provided options to the default options
func applyOptions(options []ClientOption) *clientOptions {
	opts := getDefaultOptions()
	for _, option := range options {
		option(opts)
	}
	return opts
}
