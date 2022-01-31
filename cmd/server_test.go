package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestBlockHosts(t *testing.T) {

	// Create a test backend that wraps our proxyHandler. This test backend
	// can then be sent various HTTP requests in test cases
	backend := httptest.NewServer(http.HandlerFunc(blockListAwareHandler))
	defer backend.Close()

	proxyUrl, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Create a client that will use our procrastiproxy as a proxy, so that we can
	// test our proxy's functionality
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	type TestCase struct {
		BlockedHosts       []string
		URL                string
		WantHTTPStatusCode int
	}

	testCases := []TestCase{
		{BlockedHosts: []string{"reddit.com"}, URL: "http://reddit.com", WantHTTPStatusCode: http.StatusForbidden},
		{BlockedHosts: []string{"nytimes.com"}, URL: "http://wikipedia.org", WantHTTPStatusCode: http.StatusOK},
	}

	for _, tc := range testCases {

		// Set the blocked hosts config variable that is used by the proxy backend
		viper.Set("blocked_hosts", tc.BlockedHosts)
		LoadBlockList()

		res, err := client.Get(tc.URL)
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != tc.WantHTTPStatusCode {
			t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", tc.WantHTTPStatusCode, tc.URL, res.StatusCode)
			t.Fail()
		}
	}
}
