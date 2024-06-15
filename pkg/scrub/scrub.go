// This package provides functions for hiding sensitive info.
package scrub

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/williampsena/bugs-channel-plugins/pkg/event"
)

// This function iterates event extra and stack trace recursively, hiding sensitive information and producing a safe map.
func ScrubSensitiveEvent(e *event.Event, sensitiveKeys []string) {
	e.Extra = scrubSensitiveFromMap(&e.Extra, sensitiveKeys)
	e.Tags = scrubSensitiveFromTags(&e.Tags, sensitiveKeys)
}

// This function iterates a map recursively, hiding sensitive information and producing a safe map.
func scrubSensitiveFromMap(data *map[string]interface{}, sensitiveKeys []string) map[string]interface{} {
	output := make(map[string]interface{})

	for key, value := range *data {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Array, reflect.Slice:
			if values, ok := value.([]map[string]interface{}); ok {
				value = scrubSensitiveFromArray(&values, sensitiveKeys)
			}
		case reflect.Map:
			parsedValue := value.(map[string]interface{})
			value = scrubSensitiveFromMap(&parsedValue, sensitiveKeys)
		default:
			if slices.Contains(sensitiveKeys, key) {
				value = "*"
			}
		}

		output[key] = value
	}

	return output
}

// This function iterates a array of map recursively, hiding sensitive information and producing a safe array map.
func scrubSensitiveFromArray(values *[]map[string]interface{}, sensitiveKeys []string) []map[string]interface{} {
	var output []map[string]interface{}

	for _, value := range *values {
		output = append(output, scrubSensitiveFromMap(&value, sensitiveKeys))
	}

	return output
}

func scrubSensitiveFromTags(tags *[]string, sensitiveKeys []string) []string {
	var output []string

	for _, tag := range *tags {
		values := strings.Split(tag, ":")
		key := values[0]
		value := values[1]

		if slices.Contains(sensitiveKeys, key) {
			value = "*"
		}

		output = append(output, fmt.Sprintf("%v:%v", key, value))
	}

	return output
}
