package pluck

import (
	gloo "github.com/gloo-foo/framework"

	jsoncmd "github.com/gloo-foo/cmd-json"
)

// Field names one key of a JSON object.
type Field string

// pluckFields copies the named fields that exist in obj into a fresh object.
func pluckFields(obj map[string]jsoncmd.Value, fields []Field) map[string]jsoncmd.Value {
	result := make(map[string]jsoncmd.Value, len(fields))
	for _, field := range fields {
		if value, exists := obj[string(field)]; exists {
			result[string(field)] = value
		}
	}
	return result
}

// pluck is the Processor: keep the named fields of an object, dropping the value
// entirely when it is not an object or retains none of the fields.
func pluck(fields []Field) jsoncmd.Processor {
	return func(v jsoncmd.Value) (jsoncmd.Value, bool, error) {
		obj, ok := jsoncmd.AsMap(v)
		if !ok {
			return nil, false, nil
		}
		result := pluckFields(obj, fields)
		if len(result) == 0 {
			return nil, false, nil
		}
		return result, true, nil
	}
}

// Pluck returns a command that keeps only the named fields of each JSON object,
// emitting one compact object per line. Objects that contain none of the named
// fields are dropped, as are non-object values.
//
// Shell analogue: jq '{name, age}'
func Pluck(fields ...Field) gloo.Command[[]byte, []byte] {
	return jsoncmd.Process(pluck(fields))
}
