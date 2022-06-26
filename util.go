package procrastiproxy

import (
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

func validateBlockListInput(blockListMembers []string) error {
	if len(blockListMembers) < 1 {
		return EmptyBlockListError{}
	}
	return nil
}

func parseBlockListInput(blockList *string, list *List) error {
	var blockListMembers []string
	var blockListString = *blockList
	if blockListString != "" {
		blockListMembers = strings.Split(blockListString, ",")
	}

	validationErr := validateBlockListInput(blockListMembers)
	if validationErr != nil {
		return validationErr
	}

	for _, host := range blockListMembers {
		AddHostToBlockList(list, host)
	}

	return nil
}

type timeFlag struct {
	Name  string
	Value string
}

func parseStartAndEndTimes(blockTimeStart, blockTimeEnd string) error {

	var result *multierror.Error

	timeFlags := []timeFlag{
		{
			Name:  "block-time-start",
			Value: blockTimeStart,
		},
		{
			Name:  "block-time-end",
			Value: blockTimeEnd,
		},
	}

	for _, timeFlag := range timeFlags {
		if _, parseErr := time.Parse(time.Kitchen, timeFlag.Value); parseErr != nil {
			result = multierror.Append(result, InvalidTimeFormatError{FlagName: timeFlag.Name, Value: timeFlag.Value, Underlying: parseErr})
		}
	}
	return result.ErrorOrNil()
}

// Build the fast, in-memory list of blocked hosts from the configured values
func AddHostToBlockList(list *List, hosts ...string) {
	for _, host := range hosts {
		list.Add(host)
	}
}

func SlicesAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	z := sort.StringSlice(a)
	y := sort.StringSlice(b)

	z.Sort()
	y.Sort()

	for i := range z {
		if z[i] != y[i] {
			return false
		}
	}

	return true
}
