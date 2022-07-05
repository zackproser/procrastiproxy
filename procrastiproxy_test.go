package procrastiproxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunCLIWithoutRequiredInputsErrors(t *testing.T) {
	err := RunCLI()
	require.Error(t, err)
}

func TestInvalidTimeFlagsRejected(t *testing.T) {
	type TestCase struct {
		Name           string
		BlockTimeStart string
		BlockTimeEnd   string
		Want           error
	}
	testCases := []TestCase{
		{
			Name:           "Invalid BlockTimeStart rejected",
			BlockTimeStart: "IamInvalid",
			BlockTimeEnd:   "5:00PM",
			Want:           InvalidTimeFormatError{},
		},
		{
			Name:           "Invalid BlockTimeEnd rejected",
			BlockTimeStart: "8:10AM",
			BlockTimeEnd:   "45difyr8&E&FDG",
			Want:           InvalidTimeFormatError{},
		},
		{
			Name:           "Invalid BlockTimeStart and BlockTimeEnd values rejected",
			BlockTimeStart: "ThisisntValid",
			BlockTimeEnd:   "AndNeitherIsThis",
			Want:           InvalidTimeFormatError{},
		},
		{
			Name:           "Valid start and end times accepted",
			BlockTimeStart: "8:30AM",
			BlockTimeEnd:   "6:00PM",
			Want:           nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			err := parseStartAndEndTimes(tc.BlockTimeStart, tc.BlockTimeEnd)

			if tc.Want == nil && err == nil {
				return
			}

			if !errors.As(err, &tc.Want) {
				t.Logf("%s - wanted error of type %T but got %T", tc.Name, tc.Want, err)
				t.Fail()
			}
		})
	}

}

func TestGetList(t *testing.T) {
	p := NewProcrastiproxy()

	maybeList := p.GetList()
	require.NotNil(t, maybeList)
}

func TestListAddResultsInExpectedElements(t *testing.T) {

	type TestCase struct {
		Name          string
		ElementsToAdd []string
	}

	testCases := []TestCase{
		{
			Name:          "When adding 1 element",
			ElementsToAdd: []string{"one"},
		},
		{
			Name:          "When adding 3 elements",
			ElementsToAdd: []string{"thing1", "thing2", "thing3"},
		},
		{
			Name:          "When adding 4 elements",
			ElementsToAdd: []string{"one fish", "two fish", "red fish", "blue fish"},
		},
	}

	for _, tc := range testCases {
		p := NewProcrastiproxy()
		// Reset the list singleton before each test case is run
		l := p.GetList()

		t.Run(tc.Name, func(t *testing.T) {
			for _, item := range tc.ElementsToAdd {
				l.Add(item)
			}

			require.Equal(t, len(tc.ElementsToAdd), l.Length())
			require.True(t, SlicesAreEqual(tc.ElementsToAdd, l.All()))

		})
	}

}

func TestListRemoveResultsInExpectedElements(t *testing.T) {

	type TestCase struct {
		Name             string
		ElementsToAdd    []string
		ElementsToRemove []string
		Want             []string
	}

	testCases := []TestCase{
		{
			Name:             "Removing 2 elements",
			ElementsToAdd:    []string{"one", "two", "three", "four"},
			ElementsToRemove: []string{"one", "two"},
			Want:             []string{"three", "four"},
		},
	}

	for _, tc := range testCases {

		p := NewProcrastiproxy()
		l := p.GetList()

		t.Run(tc.Name, func(t *testing.T) {

			for _, item := range tc.ElementsToAdd {
				l.Add(item)
			}

			for _, item := range tc.ElementsToRemove {
				l.Remove(item)
			}

			require.Equal(t, len(tc.Want), l.Length())

		})
	}
}

// TestHostBlocking ensures that you can access a host via the block list aware handler before it is added to the block list,
// but not after it is added to the block list
func TestHostBlocking(t *testing.T) {
	t.Parallel()

	p := NewProcrastiproxy()

	// Create a test server to mimic reddit.com
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")

	}))
	defer ts.Close()

	// First, ensure that the target host can be reached before it is added to the block list
	fr := httptest.NewRequest("GET", ts.URL, strings.NewReader(""))

	fw := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(fw, fr)

	if fw.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL prior to adding host to block list: %s but got: %d\n", http.StatusOK, ts.URL, fw.Code)
	}

	u, parseErr := url.Parse(ts.URL)
	require.NoError(t, parseErr)

	AddHostToBlockList(p.GetList(), u.Host)

	// Next, ensure the same host cannot be accessed after being added to the block list

	r := httptest.NewRequest("GET", ts.URL, strings.NewReader(""))

	w := httptest.NewRecorder()

	http.HandlerFunc(p.blockListAwareHandler).ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Logf("Wanted HTTP StatusCode: %d for URL after adding host to block list: %s but got: %d\n", http.StatusForbidden, ts.URL, w.Code)
		t.Fail()
	}
}

// TestProxiedHost ensures that you can use the proxyHandler to do "pass-through" network requests
func TestProxiedHost(t *testing.T) {
	t.Parallel()

	p := NewProcrastiproxy()

	// Create a test server to mimic reddit.com
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")

	}))
	defer ts.Close()

	// First, ensure that the target host can be reached before it is added to the block list
	r := httptest.NewRequest("GET", ts.URL, strings.NewReader(""))

	w := httptest.NewRecorder()

	http.HandlerFunc(p.proxyHandler).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Logf("Wanted HTTP StatusCode: %d for URL: %s but got: %d\n", http.StatusOK, ts.URL, w.Code)
		t.Fail()
	}
}

// TestAdminHandlerBlocksHostsDynamicallt ensures that you can dynamically add a new host to the block list
// via the admin/block endpoint, that will be respected for all subsequent requests
func TestAdminHandlerBlocksHostsDynamically(t *testing.T) {
	t.Parallel()

	p := NewProcrastiproxy()

	t.Logf("p.GetList().All(): %+v\n", p.GetList().All())

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

func TestWithinBlockWindow_IsTrue(t *testing.T) {
	testCases := []struct {
		Name      string
		CheckTime string
	}{
		{
			Name:      "At start time",
			CheckTime: "9:00AM",
		},
		{
			Name:      "1 minute after start time",
			CheckTime: "9:01AM",
		},
		{
			Name:      "1 minute before end time",
			CheckTime: "4:59PM",
		},
		{
			Name:      "2 hours after start time",
			CheckTime: "11:00AM",
		},
		{
			Name:      "3 hours after start time",
			CheckTime: "12:00PM",
		},
		{
			Name:      "4 hours after start time",
			CheckTime: "1:00PM",
		},
		{
			Name:      "1 hour before end time",
			CheckTime: "4:00PM",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a new instance of Procrastiproxy
			p := NewProcrastiproxy()

			// Configure its timing window
			p.ConfigureProxyTimeSettings("9:00AM", "5:00PM")

			parsedCheckTime, parseErr := time.Parse(time.Kitchen, tc.CheckTime)
			if parseErr != nil {
				t.Logf("Error parsing check time: %s - error: %s\n", tc.CheckTime, parseErr)
			}

			// Test the WithinBlockWindow method
			got := p.WithinBlockWindow(parsedCheckTime)

			t.Logf("tc.CheckTime: %s\n", tc.CheckTime)

			if !got {
				t.Logf("Wanted: %v for WithinBlockWindow  (%v, %v), but got: %v\n", true, tc.CheckTime, p.GetProxyTimeSettings(), got)
				t.Fail()
			}
		})
	}
}

func TestWithinBlockWindow_IsFalse(t *testing.T) {
	testCases := []struct {
		Name      string
		CheckTime string
	}{
		{
			Name:      "At end time",
			CheckTime: "5:00PM",
		},
		{
			Name:      "1 hour before start time",
			CheckTime: "8:00AM",
		},
		{
			Name:      "1 minute before start time",
			CheckTime: "8:59AM",
		},
		{
			Name:      "2 hours after end time",
			CheckTime: "7:00PM",
		},
		{
			Name:      "3 hours after end time",
			CheckTime: "8:00PM",
		},
		{
			Name:      "4 hours after end time",
			CheckTime: "9:00PM",
		},
		{
			Name:      "1 minute after end time",
			CheckTime: "5:01PM",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a new instance of Procrastiproxy
			p := NewProcrastiproxy()

			// Configure its timing window
			p.ConfigureProxyTimeSettings("9:00AM", "5:00PM")

			parsedCheckTime, parseErr := time.Parse(time.Kitchen, tc.CheckTime)
			if parseErr != nil {
				t.Logf("Error parsing check time: %s - error: %s\n", tc.CheckTime, parseErr)
			}

			// Test the WithinBlockWindow method
			got := p.WithinBlockWindow(parsedCheckTime)

			t.Logf("tc.CheckTime: %s\n", tc.CheckTime)

			if got {
				t.Logf("Wanted: %v for WithinBlockWindow  (%v, %v), but got: %v\n", false, tc.CheckTime, p.GetProxyTimeSettings(), got)
				t.Fail()
			}
		})
	}
}

func TestValidateBlockListInputErrorsOnEmptyMemberSlice(t *testing.T) {

	err := validateBlockListInput([]string{})
	require.Error(t, err)

	if emptyBlockListErr, ok := err.(EmptyBlockListError); !ok {
		t.Logf("Expected error of type: %T, but got: %T", EmptyBlockListError{}, emptyBlockListErr)
		t.Fail()
	}
}

func TestParseBlockListInput(t *testing.T) {

	type TestCase struct {
		Name        string
		InputString string
		Want        []string
	}

	testCases := []TestCase{
		{
			Name:        "single block host input",
			InputString: "reddit.com",
			Want:        []string{"reddit.com"},
		},
		{
			Name:        "two block hosts input",
			InputString: "reddit.com,nytimes.com",
			Want:        []string{"reddit.com", "nytimes.com"},
		},
		{
			Name:        "three block hosts input",
			InputString: "reddit.com,nytimes.com,twitter.com",
			Want:        []string{"reddit.com", "nytimes.com", "twitter.com"},
		},
	}

	for _, tc := range testCases {

		l := NewList()

		t.Run(tc.Name, func(t *testing.T) {
			err := parseBlockListInput(&tc.InputString, l)
			require.NoError(t, err)

			// Ensure list has expected members following parsing
			for _, member := range tc.Want {
				require.True(t, l.Contains(member))
				//Also sanity check that hostIsBlocked returns true
				require.True(t, hostIsOnBlockList(member, l))
			}

		})
	}
}

func TestParseStartAndEndTimes(t *testing.T) {

	type TestCase struct {
		Name            string
		StartTimeString string
		EndTimeString   string
	}

	testCases := []TestCase{
		{
			Name:            "top of the hour",
			StartTimeString: "9:00AM",
			EndTimeString:   "5:00PM",
		},
		{
			Name:            "minutes defined",
			StartTimeString: "9:38AM",
			EndTimeString:   "5:14PM",
		},
		{
			Name:            "15 minute block window",
			StartTimeString: "9:00AM",
			EndTimeString:   "9:15AM",
		},
		{
			Name:            "1 minute block window",
			StartTimeString: "9:00AM",
			EndTimeString:   "9:01AM",
		},
		{
			Name:            "18 hour block window",
			StartTimeString: "12:00AM",
			EndTimeString:   "6:00PM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := parseStartAndEndTimes(tc.StartTimeString, tc.EndTimeString)
			require.NoError(t, err)
		})
	}
}
