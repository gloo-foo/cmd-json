package alias_test

import (
	"testing"

	"github.com/gloo-foo/cmd-json/alias"
	"github.com/gloo-foo/testable"
)

// The alias package re-exports the Json constructor under an unprefixed name. A
// mis-wired re-export (Json bound to the wrong function) compiles cleanly, so
// only behavior can prove the wiring: the re-exported Json must compact and
// key-sort each input line exactly as command.Json does.

func TestAlias_JsonCompacts(t *testing.T) {
	lines, err := testable.TestLines(alias.Json(), `{ "b" : 2 , "a" : 1 }`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"a":1,"b":2}` {
		t.Fatalf("got %q, want compact key-sorted object", lines)
	}
}

func TestAlias_JsonInvalidErrors(t *testing.T) {
	if _, err := testable.TestLines(alias.Json(), "not json\n"); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
