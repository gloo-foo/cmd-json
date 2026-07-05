package fromtoml_test

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	"github.com/gloo-foo/cmd-json/fromtoml"
)

func TestFromToml_Basic(t *testing.T) {
	in := "name = \"Alice\"\nage = 30\n"
	got, err := testable.TestLines(fromtoml.FromToml(), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{"age":30,"name":"Alice"}` {
		t.Fatalf("got %q, want one compact object", got)
	}
}

func TestFromToml_NestedTable(t *testing.T) {
	in := "[server]\nhost = \"localhost\"\nport = 8080\n"
	got, err := testable.TestLines(fromtoml.FromToml(), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{"server":{"host":"localhost","port":8080}}` {
		t.Fatalf("got %q", got)
	}
}

func TestFromToml_EmptyInput(t *testing.T) {
	got, err := testable.TestLines(fromtoml.FromToml(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{}` {
		t.Fatalf("got %q, want {}", got)
	}
}

func TestFromToml_InvalidErrors(t *testing.T) {
	_, err := testable.TestLines(fromtoml.FromToml(), "this is = not = valid toml\n")
	if !errors.Is(err, fromtoml.ErrTOML) {
		t.Fatalf("got %v, want ErrTOML", err)
	}
}
