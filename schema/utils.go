package schema

func StringInSlice(value string, slice []string) bool {
	for _, sliceValue := range slice {
		if value == sliceValue {
			return true
		}
	}
	return false
}
