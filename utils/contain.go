package utils

func ContainsInSlice(items []bool, item bool) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
