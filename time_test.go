package procrastiproxy

import (
	"testing"
	"time"
)

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
