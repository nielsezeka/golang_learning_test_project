package utils

import (
	"fmt"

	"github.com/lib/pq"
)

// BuildUpdateQuery builds SQL set clauses and args for dynamic updates from a map.
// allowedFields maps field names to a converter function (for custom handling, e.g., pq.Array).
// Returns setClauses, args, and next arg index.
func BuildUpdateQuery(input map[string]interface{}, allowedFields map[string]func(interface{}) (interface{}, error), startIdx int) ([]string, []interface{}, int, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := startIdx

	for field, converter := range allowedFields {
		if val, ok := input[field]; ok {
			converted, err := converter(val)
			if err != nil {
				return nil, nil, 0, fmt.Errorf("invalid value for %s: %w", field, err)
			}
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIdx))
			args = append(args, converted)
			argIdx++
		}
	}
	return setClauses, args, argIdx, nil
}

// Example converter for pq.Array fields
func StringArrayConverter(val interface{}) (interface{}, error) {
	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not an array")
	}
	strArr := make([]string, len(arr))
	for i, v := range arr {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element not a string")
		}
		strArr[i] = str
	}
	return pq.Array(strArr), nil
}

// Example converter for string fields
func StringConverter(val interface{}) (interface{}, error) {
	str, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("not a string")
	}
	return str, nil
}
