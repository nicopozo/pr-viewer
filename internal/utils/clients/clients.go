package clients

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/nicopozo/pr-viewer/internal/utils/log"
)

//go:generate mockgen -destination=../utils/testutils/mocks/round_tripper_mock.go -package=mocks net/http RoundTripper
//go:generate mockgen -destination=../utils/testutils/mocks/io_read_closer_mock.go -package=mocks io ReadCloser

const (
	defaultMaxIdleConns     = 100
	defaultIdleConnTimeout  = 90 * time.Second
	defaultHandshakeTimeout = 10 * time.Second
	defaultContinueTimeout  = 1 * time.Second
	defaultDialTimeout      = 30 * time.Second

	HeaderContentTypeKey    = "Content-Type"
	HeaderAcceptKey         = "Accept"
	HeaderContentTypeJSON   = "application/json"
	HeaderContentTypeXML    = "application/xml"
	HeaderTrackingIDKey     = "X-Tracking-Id"
	HeaderCallerScopesKey   = "X-Caller-Scopes"
	HeaderCallerScopesAdmin = "admin"
	HeaderAuthorization     = "Authorization"

	DefaultAppName    = "unknown"
	DefaultAppEnv     = "unknown"
	DefaultClientName = "unknown"
)

var AdminCallerHeaders = map[string]string{ //nolint:gochecknoglobals
	HeaderCallerScopesKey: HeaderCallerScopesAdmin,
}

type RetrySettings struct {
	Delay    time.Duration
	Attempts int
}

// HTTPSettings are the settings needed to create HTTP clients.
type HTTPSettings struct {
	MaxIdleConns     int
	IdleConnTimeout  time.Duration
	HandshakeTimeout time.Duration
	ContinueTimeout  time.Duration
	DialTimeout      time.Duration
}

// NewHTTPSettings creates a default HTTP client settings.
func NewHTTPSettings() *HTTPSettings {
	return &HTTPSettings{
		MaxIdleConns:     defaultMaxIdleConns,
		IdleConnTimeout:  defaultIdleConnTimeout,
		HandshakeTimeout: defaultHandshakeTimeout,
		ContinueTimeout:  defaultContinueTimeout,
		DialTimeout:      defaultDialTimeout,
	}
}

// NewHTTPSettings creates an HTTP client with the given settings.
func NewHTTPClient(settings *HTTPSettings) http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   settings.DialTimeout,
			KeepAlive: settings.DialTimeout,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          settings.MaxIdleConns,
		IdleConnTimeout:       settings.IdleConnTimeout,
		TLSHandshakeTimeout:   settings.HandshakeTimeout,
		ExpectContinueTimeout: settings.ContinueTimeout,
	}

	return NewHTTPClientWithTransport(transport)
}

// NewHTTPClientWithTransport creates an HTTP client with the given transport.
func NewHTTPClientWithTransport(transport http.RoundTripper) http.Client {
	return http.Client{
		Timeout:   time.Minute,
		Transport: transport,
	}
}

func CloseResponseBody(cli interface{}, response *http.Response, logger log.ILogger) {
	if response != nil && response.Body != nil {
		if err := response.Body.Close(); err != nil {
			errorMsg := "error closing response body"

			if logger != nil {
				logger.Error(cli, nil, err, errorMsg)
			} else {
				fmt.Printf("%T %s\n", cli, errorMsg) //nolint:forbidigo
			}
		}
	}
}

func NewRequest(ctx context.Context, httpMethod, url string, body io.Reader, params,
	headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, body)

	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	addParameters(req, params)
	addJSONHeaders(req)
	addHeaders(req, headers)

	req.Close = true

	return req, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	if len(headers) > 0 {
		for headerName, headerValue := range headers {
			req.Header.Set(headerName, headerValue)
		}
	}
}

func addParameters(req *http.Request, params map[string]string) {
	if len(params) > 0 {
		rawQuery := req.URL.Query()

		for paramKey, paramValue := range params {
			rawQuery.Add(paramKey, paramValue)
		}

		req.URL.RawQuery = rawQuery.Encode()
	}
}

func addJSONHeaders(req *http.Request) {
	req.Header.Set(HeaderContentTypeKey, HeaderContentTypeJSON)
	req.Header.Set(HeaderAcceptKey, HeaderContentTypeJSON)
}
