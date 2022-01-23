package cmd

func includes(haystack []string, needle string) bool {
	for _, member := range haystack {
		if member == needle {
			return true
		}
	}
	return false
}
