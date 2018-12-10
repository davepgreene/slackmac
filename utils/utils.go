package utils

import (
	"reflect"
	"runtime"
)

// GetFunctionName uses reflection to get the name of a function as a string
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// GetTypeName uses reflection to get the name of a type as a string
func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

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
