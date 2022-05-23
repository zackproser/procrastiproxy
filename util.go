package procrastiproxy

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.Info("Error parsing time string:", err)
	}
	return parsed
}
