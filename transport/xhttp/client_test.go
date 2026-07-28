package xhttp

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/metacubex/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestPacketUpWriterGetMethod(t *testing.T) {
	config := Config{
		Host:                "example.com",
		Path:                "/xhttp",
		UplinkHTTPMethod:    http.MethodGet,
		UplinkDataPlacement: PlacementHeader,
		UplinkDataKey:       "X-Data",
		UplinkChunkSize:     "3000",
		XPaddingBytes:       "0",
	}
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, "/xhttp/session/0", request.URL.Path)
		assert.Empty(t, request.Body)
		assert.NotEmpty(t, request.Header.Get("X-Data-0"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewReader(nil)),
		}, nil
	})
	writer := PacketUpWriter{
		ctx:       context.Background(),
		cfg:       &config,
		sessionID: "session",
		transport: transport,
	}

	n, err := writer.write([]byte("payload"))

	require.NoError(t, err)
	assert.Equal(t, len("payload"), n)
}
