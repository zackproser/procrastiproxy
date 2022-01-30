package cmd

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
		StartBlockTime string
		EndBlockTime   string
		CheckTime      time.Time
		Want           bool
	}
	testCases := []TestCase{
		{
			StartBlockTime: "9:00AM",
			EndBlockTime:   "5:00PM",
			CheckTime:      parseTime("10:00AM"),
			Want:           true,
		},
	}
	for _, tc := range testCases {
		pts := ProxyTimeSettings{
			BlockStartTime: tc.StartBlockTime,
			BlockEndTime:   tc.EndBlockTime,
		}
		got := WithinBlockWindow(tc.CheckTime, pts)
		if got != tc.Want {
			t.Logf("Wanted: %v for WithinBlockWindow(%v, %v), but got: %v\n", tc.Want, tc.CheckTime, pts, got)
			t.Fail()
		}
	}
}
