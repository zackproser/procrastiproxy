package procrastiproxy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

		l := GetList()
		l.Clear()

		t.Run(tc.Name, func(t *testing.T) {
			err := parseBlockListInput(&tc.InputString)
			require.NoError(t, err)

			// Ensure list has expected members following parsing
			l := GetList()
			for _, member := range tc.Want {
				require.True(t, l.Contains(member))
				//Also sanity check that hostIsBlocked returns true
				require.True(t, hostIsBlocked(member))
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
