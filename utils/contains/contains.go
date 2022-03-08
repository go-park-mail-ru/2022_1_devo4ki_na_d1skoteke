package contains

func Contains(slice []string, str string) bool {
	for _, elem := range slice {
		if str == elem {
			return true
		}
	}
	return false
}
