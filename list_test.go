package procrastiproxy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetList(t *testing.T) {
	maybeList := GetList()
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
		// Reset the list singleton before each test case is run
		l := GetList()
		l.Clear()
		t.Run(tc.Name, func(t *testing.T) {

			// Instantiate a new list, then add every test element to it sequentially
			l := GetList()
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

		l := GetList()
		l.Clear()

		t.Run(tc.Name, func(t *testing.T) {

			// Instantiate a new list, then add every test element to it sequentially
			l := GetList()
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
