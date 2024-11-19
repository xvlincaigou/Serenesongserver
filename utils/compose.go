package utils

import (
	// "github.com/gin-gonic/gin"
	// "log"
	// "reflect"
)

func ToStringArray(value interface{}) []string {
	if value == nil {
		return []string{}
	}
	if slice, ok := value.([]string); ok {
		return slice
	}
	if slice, ok := value.([]interface{}); ok {
		result := make([]string, len(slice))
		for i, v := range slice {
			result[i] = v.(string)
		}
		return result
	}
	return []string{}
}