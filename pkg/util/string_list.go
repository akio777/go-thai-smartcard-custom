package util

func StringExistsInSlice(target string, slice []string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}
