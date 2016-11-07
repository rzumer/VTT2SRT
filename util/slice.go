package util

func Contains(slice []byte, item byte) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	
	return false
}
