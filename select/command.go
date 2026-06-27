package selectcmd

import (
	jsoncmd "github.com/gloo-foo/cmd-json"
	gloo "github.com/gloo-foo/framework"
)

// Condition reports whether a decoded JSON value should be selected.
type Condition func(value jsoncmd.Value) bool

// Select returns a command that keeps each JSON value for which condition
// reports true, emitting one compact value per line. Values that fail the
// condition are dropped.
//
// Shell analogue: jq 'select(...)'
func Select(condition Condition) gloo.Command[[]byte, []byte] {
	return jsoncmd.Process(func(v jsoncmd.Value) (jsoncmd.Value, bool, error) {
		return v, condition(v), nil
	})
}

// Common condition builders

// HasField returns a condition that checks if an object has a specific field
func HasField(field string) Condition {
	return func(value any) bool {
		obj, ok := value.(map[string]any)
		if !ok {
			return false
		}
		_, exists := obj[field]
		return exists
	}
}

// FieldEquals returns a condition that checks if a field equals a specific value
func FieldEquals(field string, expected any) Condition {
	return func(value any) bool {
		obj, ok := value.(map[string]any)
		if !ok {
			return false
		}
		actual, exists := obj[field]
		return exists && actual == expected
	}
}

// FieldMatches returns a condition that checks if a field matches a predicate
func FieldMatches(field string, predicate func(any) bool) Condition {
	return func(value any) bool {
		obj, ok := value.(map[string]any)
		if !ok {
			return false
		}
		fieldValue, exists := obj[field]
		return exists && predicate(fieldValue)
	}
}

// And combines multiple conditions with logical AND
func And(conditions ...Condition) Condition {
	return func(value any) bool {
		for _, cond := range conditions {
			if !cond(value) {
				return false
			}
		}
		return true
	}
}

// Or combines multiple conditions with logical OR
func Or(conditions ...Condition) Condition {
	return func(value any) bool {
		for _, cond := range conditions {
			if cond(value) {
				return true
			}
		}
		return false
	}
}

// Not negates a condition
func Not(condition Condition) Condition {
	return func(value any) bool {
		return !condition(value)
	}
}
