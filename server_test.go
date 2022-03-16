package procrastiproxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setupTestServer is a convenience method for creating a test backend that we can make
// HTTP requests to
func setupTestServer(t *testing.T) (*http.Client, *httptest.Server, error) {
	// Create a test backend that wraps our proxyHandler. This test backend
	// can then be sent various HTTP requests in test cases
	backend := httptest.NewServer(http.HandlerFunc(blockListAwareHandler))

	proxyURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Create a client that will use our procrastiproxy as a proxy, so that we can
	// test our proxy's functionality
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	return client, backend, nil
}

// TestDeniedBlockHosts adds a host to the block list and then immediately attempts to make
// a request to that host, which should be denied by procrastiproxy
func TestDeniesBlockedHost(t *testing.T) {
	client, backend, err := setupTestServer(t)
	defer backend.Close()
	if err != nil {
		t.Fatal(err)
	}

	blockedHost := "reddit.com"
	testURL := "http://reddit.com"

	AddHostToBlockList(blockedHost)

	res, err := client.Get(testURL)
	if err != nil {
		t.Fatalf("Error attempting to fetch test server URL: %v\n", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusForbidden, testURL, res.StatusCode)
		t.Fail()
	}
}

// TestAllowsWhitelistedHost ensures that a host that has not been explicitly blocked
// can be reached through procrastiproxy
func TestAllowsWhitelistedHost(t *testing.T) {

	client, backend, err := setupTestServer(t)
	defer backend.Close()
	if err != nil {
		t.Fatal(err)
	}

	blockedHost := "nytimes.com"
	testURL := "http://wikipedia.org"

	// Set the blocked hosts config variable that is used by the proxy backend
	AddHostToBlockList(blockedHost)

	res, err := client.Get(testURL)
	if err != nil {
		t.Fatalf("Error attempting to fetch test server URL: %v\n", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testURL, res.StatusCode)
		t.Fail()
	}
}
