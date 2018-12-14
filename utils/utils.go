package utils

// MapKeys generates a string slice of a map[string]T's keys
func MapKeys(i interface{}) []string {
	keys := make([]string, 0)
	switch x := i.(type) {
	case map[string]interface{}:
		for k := range x {
			keys = append(keys, k)
		}
	default:
		//
	}

	return keys
}
