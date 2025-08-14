package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	serverURL      = "http://localhost:8089"
	nginxURL       = "http://nginx:80"
	defaultTimeout = 10 * time.Second
)

func doGet(t *testing.T, timeout time.Duration, url string) *http.Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err, "failed to create request for %s", url)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "request to %s failed", url)

	return resp
}

func TestImageFoundInCache(t *testing.T) {
	url := fmt.Sprintf("%s/fill/250/250/%s/images/_gopher_original_1024x504.jpg", serverURL, nginxURL)

	resp := doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200")
	require.Equal(t, "MISS", resp.Header.Get("X-Cache"), "expected cache MISS")

	resp = doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200")
	require.Equal(t, "HIT", resp.Header.Get("X-Cache"), "expected cache HIT")
}

func TestRemoteServerNotFound(t *testing.T) {
	url := fmt.Sprintf("%s/fill/300/300/http://nonexistent:80/images/test.jpg", serverURL)

	resp := doGet(t, 5*time.Second, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "expected 502")
}

func TestImageNotFound(t *testing.T) {
	url := fmt.Sprintf("%s/fill/300/300/%s/images/nonexistent.jpg", serverURL, nginxURL)

	resp := doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "expected 502 for missing image")
}

func TestNotImage(t *testing.T) {
	url := fmt.Sprintf("%s/fill/300/300/%s/images/not-an-image", serverURL, nginxURL)

	resp := doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "expected 502 for non-image file")
}

func TestRemoteServerError(t *testing.T) {
	url := fmt.Sprintf("%s/fill/300/300/%s/images/error.jpg", serverURL, nginxURL)

	resp := doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "expected 502 for server error")
}

func TestImageOK(t *testing.T) {
	url := fmt.Sprintf("%s/fill/120/300/%s/images/_gopher_original_1024x504.jpg", serverURL, nginxURL)

	resp := doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200")
	require.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"), "expected Content-Type image/jpeg")
}

func TestInvalidParams(t *testing.T) {
	url := fmt.Sprintf("%s/fill/invalid/300/%s/images/_gopher_original_1024x504.jpg", serverURL, nginxURL)

	resp := doGet(t, defaultTimeout, url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected 400 for invalid params")
}
