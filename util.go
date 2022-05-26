package procrastiproxy

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-multierror"
)

func validateBlockListInput(blockListMembers []string) error {
	if len(blockListMembers) < 1 {
		return EmptyBlockListError{}
	}
	return nil
}

func parseBlockListInput(blockList *string) error {
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
		AddHostToBlockList(host)
	}

	return nil
}

type timeFlag struct {
	Name  string
	Value string
}

func parseStartAndStopTimes(blockTimeStart, blockTimeEnd string) error {

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
func AddHostToBlockList(hosts ...string) {
	list := GetList()
	for _, host := range hosts {
		list.Add(host)
	}
}

func GetBlockedHosts() []string {
	return GetList().All()
}

func includes(haystack []string, needle string) bool {
	for _, member := range haystack {
		if member == needle {
			return true
		}
	}
	return false
}

func parseTime(timeString string) time.Time {
	parsed, err := time.Parse(time.Kitchen, timeString)
	if err != nil {
		log.Debug("Error parsing time string:", err)
	}
	return parsed
}
