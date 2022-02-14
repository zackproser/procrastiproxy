package procrastiproxy

import "strings"

func parseBlockListInput(blockList *string) {
	var blockListMembers []string
	var blockListString = *blockList
	if blockListString != "" {
		blockListMembers = strings.Split(blockListString, ",")
	}
	for _, host := range blockListMembers {
		AddHostToBlockList(host)
	}
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
