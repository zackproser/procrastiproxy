package procrastiproxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T, handlerFunc func(http.ResponseWriter, *http.Request)) (*http.Client, *httptest.Server, error) {
	// Create a test backend that wraps our blockListAwareHandler. This test backend
	// can then be sent various HTTP requests in test cases
	backend := httptest.NewServer(http.HandlerFunc(handlerFunc))

	proxyURL, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Create a client that will use our procrastiproxy as a proxy, so that we can
	// test our proxy's functionality
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}

	return client, backend, nil

}

func setupBlockListAwareServer(t *testing.T) (*http.Client, *httptest.Server, error) {
	return setupTestServer(t, NewProcrastiproxy().blockListAwareHandler)
}

func setupProxyTestServer(t *testing.T) (*http.Client, *httptest.Server, error) {
	return setupTestServer(t, proxyHandler)
}

func setupAdminTestServer(t *testing.T) (*http.Client, *httptest.Server, error) {
	return setupTestServer(t, NewProcrastiproxy().adminHandler)
}

// TestDeniedBlockHosts adds a host to the block list and then immediately attempts to make
// a request to that host, which should be denied by procrastiproxy
func TestDeniesBlockedHost(t *testing.T) {
	client, backend, err := setupBlockListAwareServer(t)
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

	client, backend, err := setupBlockListAwareServer(t)
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

func TestProxiedHost(t *testing.T) {

	client, backend, err := setupProxyTestServer(t)
	defer backend.Close()
	if err != nil {
		t.Fatal(err)
	}

	testURL := "http://reddit.com"

	res, err := client.Get(testURL)
	if err != nil {
		t.Fatalf("Error attempting to fetch proxy test server URL: %v\n", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testURL, res.StatusCode)
		t.Fail()
	}
}

func TestAdminHandlerBlocksHostsDynamically(t *testing.T) {

	client, backend, err := setupAdminTestServer(t)
	defer backend.Close()
	if err != nil {
		t.Fatal(err)
	}

	testHost := "docker.com"

	testURL := fmt.Sprintf("http://localhost:8000/admin/block/%s", testHost)

	res, err := client.Get(testURL)
	if err != nil {
		t.Fatalf("Error attempting to fetch admin test server URL: %v\n", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testURL, res.StatusCode)
		t.Fail()
	}

	// Finally, ensure the host we just dynamically added to the block list is found
	l := GetList()

	require.True(t, l.Contains(testHost))

	l.Clear()
}

func TestAdminHandlerUnblocksHostsDynamically(t *testing.T) {

	client, backend, err := setupAdminTestServer(t)
	defer backend.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Start off by pre-populating the list with the test host
	testHost := "wikipedia.com"
	l := GetList()
	l.Add(testHost)

	testURL := fmt.Sprintf("http://localhost:8000/admin/unblock/%s", testHost)

	res, err := client.Get(testURL)
	if err != nil {
		t.Fatalf("Error attempting to fetch admin test server URL: %v\n", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testURL, res.StatusCode)
		t.Fail()
	}

	// Finally, ensure the host we just dynamically added to the block list is found
	require.False(t, l.Contains(testHost))

	l.Clear()
}
