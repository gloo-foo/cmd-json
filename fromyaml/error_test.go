package fromyaml

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

func TestError_Error(t *testing.T) {
	if errYAML.Error() != "yaml" {
		t.Fatalf("got %q, want %q", errYAML.Error(), "yaml")
	}
}

func TestFromYaml_InvalidYAMLErrors(t *testing.T) {
	// A bare unterminated flow mapping is not valid YAML.
	_, err := testable.TestLines(FromYaml(), "{ unbalanced\n")
	if !errors.Is(err, errYAML) {
		t.Fatalf("got %v, want errYAML", err)
	}
}

func TestFromYaml_NonStringKeysErrorOnJSON(t *testing.T) {
	// YAML allows integer mapping keys; JSON does not. Decoding succeeds but
	// re-encoding to JSON must fail, exercising the json error branch.
	_, err := testable.TestLines(FromYaml(), "1: a\n2: b\n")
	if !errors.Is(err, errJSON) {
		t.Fatalf("got %v, want errJSON", err)
	}
}
