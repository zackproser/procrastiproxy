package procrastiproxy

import (
	"errors"
	"testing"

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
