package multipart

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	pathItemOpen  = `[`
	pathItemClose = `]`
)

// toJSON converts string value into valid json value.
func toJSON(strval string) (interface{}, error) {
	// get actual values
	var jsonValue interface{}
	// if unable to unmarshal as is
	if asIsErr := json.Unmarshal([]byte(strval), &jsonValue); asIsErr != nil {
		// try to decode as a string
		if strErr := json.Unmarshal([]byte(`"`+strval+`"`), &jsonValue); strErr != nil {
			// when we unable to convert even to a string - return Unmarshal error
			// it can happen when "strval" contains special symbols, quotes etc.
			return nil, strErr
		}
	}
	return jsonValue, nil
}

// isArrayElement checks if provided  path element can be converted to an array index.
// Otherwise it is an object key.
func isArrayElement(element string) (int, bool) {
	// try to get array index
	i, err := strconv.Atoi(element)
	// empty brackets "[]" convert to element with index "0"
	if err == nil || element == "" {
		return i, true
	}
	// otherwise - this is an object key
	return i, false
}

// getObjPath explodes form path/name into elements.
func getObjPath(pathStr string) ([]string, error) {
	var path []string
	// split parts by "["
	parts := strings.Split(pathStr, pathItemOpen)
	// the first element is ready to use
	path = append(path, parts[0])
	// clean the rest
	for _, idx := range parts[1:] {
		// if token does not have "]" at the end - it has invalid format
		if !strings.HasSuffix(idx, pathItemClose) {
			return nil, fmt.Errorf("invalid property name: %q", pathStr)
		}
		// add path element to a slice
		path = append(path, idx[:len(idx)-1])
	}
	return path, nil
}

// ActualType gets the actual underlying type of field value.
func ActualType(v reflect.Value) (reflect.Value, reflect.Kind) {
	switch v.Kind() {
	// pointer
	case reflect.Ptr:
		if v.IsNil() {
			return v, reflect.Ptr
		}
		return ActualType(v.Elem())
	// interface
	case reflect.Interface:
		if v.IsNil() {
			return v, reflect.Interface
		}
		return ActualType(v.Elem())
	// actual
	default:
		return v, v.Kind()
	}
}
