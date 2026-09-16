package selectcmd_test

import (
	"testing"

	selectcmd "github.com/gloo-foo/cmd-json/select"
)

// Each condition must return false for a non-object value (a scalar, here),
// covering the type-assertion guard in HasField, FieldEquals, and FieldMatches.

func TestHasField_NonObjectIsFalse(t *testing.T) {
	if selectcmd.HasField("age")("scalar") {
		t.Fatal("HasField must be false for a non-object value")
	}
}

func TestFieldEquals_NonObjectIsFalse(t *testing.T) {
	if selectcmd.FieldEquals("status", "active")(42.0) {
		t.Fatal("FieldEquals must be false for a non-object value")
	}
}

func TestFieldEquals_MissingFieldIsFalse(t *testing.T) {
	obj := map[string]any{"name": "Alice"}
	if selectcmd.FieldEquals("status", "active")(obj) {
		t.Fatal("FieldEquals must be false when the field is absent")
	}
}

func TestFieldMatches_NonObjectIsFalse(t *testing.T) {
	matchAny := func(any) bool { return true }
	if selectcmd.FieldMatches("age", matchAny)([]any{1.0, 2.0}) {
		t.Fatal("FieldMatches must be false for a non-object value")
	}
}
