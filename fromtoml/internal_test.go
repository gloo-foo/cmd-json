package fromtoml

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

// TestFromToml_MarshalError covers the JSON encoding branch. TOML decodes only
// into JSON-marshalable Go values, so a failing encoder must be injected to
// reach this path.
func TestFromToml_MarshalError(t *testing.T) {
	original := marshal
	marshal = func(any) ([]byte, error) { return nil, errJSON }
	t.Cleanup(func() { marshal = original })

	_, err := testable.TestLines(FromToml(), "name = \"Alice\"\n")
	if !errors.Is(err, errJSON) {
		t.Fatalf("got %v, want errJSON", err)
	}
}
