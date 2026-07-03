package selectcmd

import (
	jsoncmd "github.com/gloo-foo/cmd-json"
)

// Condition reports whether a decoded JSON value should be selected.
type Condition func(value jsoncmd.Value) bool

// Field names one key of a JSON object.
type Field string

// Common condition builders

// HasField returns a condition that checks if an object has a specific field
func HasField(field Field) Condition {
	return func(value any) bool {
		obj, ok := value.(map[string]any)
		if !ok {
			return false
		}
		_, exists := obj[string(field)]
		return exists
	}
}

// FieldEquals returns a condition that checks if a field equals a specific value
func FieldEquals(field Field, expected any) Condition {
	return func(value any) bool {
		obj, ok := value.(map[string]any)
		if !ok {
			return false
		}
		actual, exists := obj[string(field)]
		return exists && actual == expected
	}
}

// FieldMatches returns a condition that checks if a field matches a predicate
func FieldMatches(field Field, predicate func(any) bool) Condition {
	return func(value any) bool {
		obj, ok := value.(map[string]any)
		if !ok {
			return false
		}
		fieldValue, exists := obj[string(field)]
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
