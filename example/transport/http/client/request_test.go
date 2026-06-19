package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/neo532/gokit/crypt/marshaler/xml"
	"github.com/neo532/gokit/metadata"
	clt "github.com/neo532/gokit/transport/http/client"
)

type testRequest struct {
	Name string `xml:"name,omitempty"`
}

type testResponse struct {
	Message string `xml:"message,omitempty"`
}

func TestRequest_Do(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		assert.Equal(t, "POST", r.Method)

		// Verify request body
		var req testRequest
		err := xml.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, "test", req.Name)

		// Return response
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		xml.NewEncoder(w).Encode(testResponse{Message: "success"})
	}))
	defer server.Close()

	// Create client and request
	oClient := clt.NewClient()
	req := clt.NewRequest(oClient,
		clt.WithUrl(server.URL),
		clt.WithMethod("POST"),
		clt.WithContentType("application/xml"),
	)

	// Execute request
	ctx := context.Background()
	request := testRequest{Name: "test"}
	response := &testResponse{}

	ctx, err := req.Do(ctx, request, response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Message)
}

func TestRequest_Do_WithRetry(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		xml.NewEncoder(w).Encode(testResponse{Message: "success"})
	}))
	defer server.Close()

	oClient := clt.NewClient()
	req := clt.NewRequest(oClient,
		clt.WithUrl(server.URL),
		clt.WithMethod("POST"),
		clt.WithContentType("application/xml"),
		clt.WithRetryTimes(1),
		clt.WithRetryDuration(time.Millisecond*100),
	)

	ctx := context.Background()
	request := testRequest{Name: "test"}
	response := &testResponse{}

	ctx, err := req.Do(ctx, request, response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Message)
	assert.Equal(t, 2, retryCount)

	fmt.Println(metadata.FromClientResponseContext(ctx))
}

func TestRequest_Do_WithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	oClient := clt.NewClient()
	req := clt.NewRequest(oClient,
		clt.WithUrl(server.URL),
		clt.WithMethod("POST"),
		clt.WithContentType("application/xml"),
	)

	ctx := context.Background()
	request := testRequest{Name: "test"}
	response := &testResponse{}

	ctx, err := req.Do(ctx, request, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "400 Bad Request")
}

func TestRequest_Do_WithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	oClient := clt.NewClient()
	req := clt.NewRequest(oClient,
		clt.WithUrl(server.URL),
		clt.WithMethod("POST"),
		clt.WithContentType("application/xml"),
		clt.WithTimeLimit(time.Second),
	)

	ctx := context.Background()
	request := testRequest{Name: "test"}
	response := &testResponse{}

	ctx, err := req.Do(ctx, request, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
