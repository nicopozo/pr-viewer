package clients_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mercadolibre/fury_recon-clients-go/pkg/clients"
	"github.com/mercadolibre/fury_recon-clients-go/pkg/utils/testutils/mocks"
	"github.com/mercadolibre/fury_recon-commons-go/pkg/log"
	"github.com/mercadolibre/fury_recon-test-utils-go/pkg/assert"
)

func TestNewHTTPClient(t *testing.T) {
	tests := []struct {
		name     string
		settings *clients.HTTPSettings
	}{
		{
			name:     "Client created successfully",
			settings: clients.NewHTTPSettings(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clients.NewHTTPClient(tt.settings)

			assert.Equals(t, "Client Timeout is not the expected", time.Minute, got.Timeout)
		})
	}
}

func TestNewHTTPSettings(t *testing.T) {
	tests := []struct {
		name                 string
		wantMaxIdleConns     int
		wantIdleConnTimeout  time.Duration
		wantHandshakeTimeout time.Duration
		wantContinueTimeout  time.Duration
		wantDialTimeout      time.Duration
	}{
		{
			name:                 "HTTP settings created successfully",
			wantMaxIdleConns:     100,
			wantIdleConnTimeout:  90 * time.Second,
			wantHandshakeTimeout: 10 * time.Second,
			wantContinueTimeout:  1 * time.Second,
			wantDialTimeout:      30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clients.NewHTTPSettings()

			assert.Equals(t, "NewHTTPSettings() MaxIdleConns is not the expected", tt.wantMaxIdleConns, got.MaxIdleConns)             //nolint
			assert.Equals(t, "NewHTTPSettings() IdleConnTimeout is not the expected", tt.wantIdleConnTimeout, got.IdleConnTimeout)    //nolint
			assert.Equals(t, "NewHTTPSettings() HandshakeTimeout is not the expected", tt.wantHandshakeTimeout, got.HandshakeTimeout) //nolint
			assert.Equals(t, "NewHTTPSettings() ContinueTimeout is not the expected", tt.wantContinueTimeout, got.ContinueTimeout)    //nolint
			assert.Equals(t, "NewHTTPSettings() DialTimeout is not the expected", tt.wantDialTimeout, got.DialTimeout)                //nolint
		})
	}
}

func TestCloseResponseBody(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		err      error
	}{
		{
			name:     "CloseResponseBody() close body successfully",
			response: &http.Response{},
			err:      nil,
		},
		{
			name:     "CloseResponseBody() close body with error",
			response: &http.Response{},
			err:      errors.New("the closing body error"),
		},
		{
			name:     "CloseResponseBody() avoid close body with nil http response",
			response: nil,
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			closerMocker := mocks.NewMockReadCloser(mockCtrl)
			defer mockCtrl.Finish()

			if tt.response != nil {
				closerMocker.EXPECT().Close().Return(tt.err)
				tt.response.Body = closerMocker
			}

			clients.CloseResponseBody(gomock.Any(), tt.response, log.DefaultLogger())
		})
	}
}
