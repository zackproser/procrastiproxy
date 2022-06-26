package procrastiproxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHostBlocking ensures that you can access a host via the block list aware handler before it is added to the block list,
// but not after it is added to the block list
func TestHostBlocking(t *testing.T) {
	t.Parallel()

	p := NewProcrastiproxy()

	blockedHost := "reddit.com"
	testURL := "http://reddit.com"

	// First, ensure that the target host can be reached before it is added to the block list
	fr := httptest.NewRequest("GET", testURL, strings.NewReader(""))

	fw := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(fw, fr)

	if fw.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL prior to adding host to block list: %s but got: %d\n", http.StatusOK, testURL, fw.Code)
	}

	AddHostToBlockList(p.GetList(), blockedHost)

	// Next, ensure the same host cannot be accessed after being added to the block list

	r := httptest.NewRequest("GET", testURL, strings.NewReader(""))

	w := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Logf("Wanted HTTP StatusCode: %d for URL after adding host to block list: %s but got: %d\n", http.StatusForbidden, testURL, w.Code)
		t.Fail()
	}
}

// TestProxiedHost ensures that you can use the proxyHandler to do "pass-through" network requests
func TestProxiedHost(t *testing.T) {
	t.Parallel()

	testURL := "http://reddit.com"

	// First, ensure that the target host can be reached before it is added to the block list
	r := httptest.NewRequest("GET", testURL, strings.NewReader(""))

	w := httptest.NewRecorder()

	http.HandlerFunc(proxyHandler).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testURL, w.Code)
		t.Fail()
	}
}

// TestAdminHandlerBlocksHostsDynamicallt ensures that you can dynamically add a new host to the block list
// via the admin/block endpoint, that will be respected for all subsequent requests
func TestAdminHandlerBlocksHostsDynamically(t *testing.T) {
	t.Parallel()

	p := NewProcrastiproxy()

	// Sanity check the initial block list is empty
	require.Equal(t, p.GetList().Length(), 0)

	testHost := "docker.com"

	// Sanity check that we can initially reach the target host because it has not yet been blocked
	testHostURL := "http://docker.com"
	fr := httptest.NewRequest("GET", testHostURL, strings.NewReader(""))
	rw := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(rw, fr)

	if rw.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testHostURL, rw.Code)
		t.Fail()
	}

	// Now, dynamically block the same target host by making a request to the admin block endpoint, passing the host info

	testURL := fmt.Sprintf("http://localhost:8000/admin/block/%s", testHost)

	r := httptest.NewRequest("GET", testURL, strings.NewReader(""))
	w := httptest.NewRecorder()

	http.HandlerFunc(p.adminHandler).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testURL, w.Code)
		t.Fail()
	}

	// Finally, ensure the host we just dynamically added to the block list is found
	require.True(t, p.GetList().Contains(testHost))
	require.Equal(t, p.GetList().Length(), 1)

	// Ensure that, attempting to hit the same host now fails because it is blocked at the proxy level
	ar := httptest.NewRequest("GET", testHostURL, strings.NewReader(""))
	aw := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(aw, ar)

	if aw.Code != http.StatusForbidden {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusForbidden, testHostURL, aw.Code)
		t.Fail()
	}
}

// TestAdminHandlerUnblocksHostsDynamically ensures that a blocked host can be removed from the block list via the admin/unblock endpoint
func TestAdminHandlerUnblocksHostsDynamically(t *testing.T) {
	t.Parallel()

	p := NewProcrastiproxy()

	testHost := "wikipedia.com"
	testHostURL := fmt.Sprintf("http://%s", testHost)

	// Pre-populate the list with the test host
	p.GetList().Add(testHost)

	require.True(t, p.GetList().Contains(testHost))

	// Sanity check that you cannot get the blocked test host to begin with
	fr := httptest.NewRequest("GET", testHostURL, strings.NewReader(""))
	rw := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(rw, fr)

	if rw.Code != http.StatusForbidden {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusForbidden, testHostURL, rw.Code)
		t.Fail()
	}

	// Now, dynamically unblock the same test host by making a request to the admin unblock endpoint
	testAdminURL := fmt.Sprintf("http://localhost:8000/admin/unblock/%s", testHost)

	r := httptest.NewRequest("GET", testAdminURL, strings.NewReader(""))
	w := httptest.NewRecorder()

	http.HandlerFunc(p.adminHandler).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testAdminURL, w.Code)
		t.Fail()
	}

	// Ensure the host we just dynamically removed from the block list is no longer found on it
	require.False(t, p.GetList().Contains(testHost))
	require.Equal(t, p.GetList().Length(), 0)

	// Ensure that we can access the test host now that it has been unblocked via the admin endpoint
	ar := httptest.NewRequest("GET", testHostURL, strings.NewReader(""))
	aw := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(aw, ar)

	if aw.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, testHostURL, aw.Code)
		t.Fail()
	}
}
