package procrastiproxy

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func parseTime(timeString string) time.Time {
	parsed, err := time.Parse(time.Kitchen, timeString)
	if err != nil {
		log.Info("Error parsing time string:", err)
	}
	return parsed
}

func TestWithinBlockWindow(t *testing.T) {
	type TestCase struct {
		Name           string
		StartBlockTime string
		EndBlockTime   string
		CheckTime      time.Time
		Want           bool
	}
	testCases := []TestCase{
		{
			Name:           "1 hour inside block window results in block",
			StartBlockTime: "9:00AM",
			EndBlockTime:   "5:00PM",
			CheckTime:      parseTime("10:00AM"),
			Want:           true,
		},
		{
			Name:           "1 hour before end of block window results in block",
			StartBlockTime: "1:00PM",
			EndBlockTime:   "5:00PM",
			CheckTime:      parseTime("4:00PM"),
			Want:           true,
		},
		{
			Name:           "1 minute after block window results in block",
			StartBlockTime: "8:00AM",
			EndBlockTime:   "5:00PM",
			CheckTime:      parseTime("8:01AM"),
			Want:           true,
		},
		{
			Name:           "1 hour after block window results in no block",
			StartBlockTime: "8:00AM",
			EndBlockTime:   "5:00PM",
			CheckTime:      parseTime("6:00PM"),
			Want:           false,
		},
		{
			Name:           "2 hours after block window results in now block",
			StartBlockTime: "8:00AM",
			EndBlockTime:   "5:00PM",
			CheckTime:      parseTime("7:00AM"),
			Want:           false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Create a new instance of Procrastiproxy
			p := NewProcrastiproxy()

			// Configure its timing window
			p.ConfigureProxyTimeSettings(tc.StartBlockTime, tc.EndBlockTime)

			// Override the default value Procrastiproxy would use for the Now() function,
			// so that we can specify a static CheckTime in our test cases
			p.Now = func() time.Time {
				return tc.CheckTime
			}

			// Test the WithinBlockWindow method
			got := p.WithinBlockWindow()
			if got != tc.Want {
				t.Logf("Wanted: %v for WithinBlockWindow  (%v, %v), but got: %v\n", tc.Want, tc.CheckTime, p.GetProxyTimeSettings(), got)
				t.Fail()
			}
		})
	}
}
