package expr

import "strings"

func convertFloat(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	} else if v, ok := value.(string); ok {
		return float64(len(v))
	} else if v, ok := value.(bool); ok && v {
		return 1
	}

	return 0
}

func convertBool(value interface{}) bool {
	if v, ok := value.(float64); ok {
		return v > 0
	} else if v, ok := value.(string); ok {
		return strings.ToLower(v) == "true"
	} else if v, ok := value.(bool); ok && v {
		return v
	}

	return false
}
