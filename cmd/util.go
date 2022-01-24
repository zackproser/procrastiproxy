package cmd

import "github.com/spf13/viper"

// Build the fast, in-memory list of blocked hosts from the configured values
func LoadBlockList() {
	list := GetList()
	for _, host := range viper.GetStringSlice("blocked_hosts") {
		list.Add(host)
	}
}

func includes(haystack []string, needle string) bool {
	for _, member := range haystack {
		if member == needle {
			return true
		}
	}
	return false
}
