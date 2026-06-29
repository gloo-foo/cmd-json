package fromyaml_test

import (
	"testing"

	"github.com/gloo-foo/testable"

	"github.com/gloo-foo/cmd-json/fromyaml"
)

func TestFromYaml_Object(t *testing.T) {
	in := "name: Alice\nage: 30\n"
	got, err := testable.TestLines(fromyaml.FromYaml(), in)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{"age":30,"name":"Alice"}` {
		t.Fatalf("got %q, want one compact object", got)
	}
}

func TestFromYaml_Nested(t *testing.T) {
	in := "server:\n  host: localhost\n  port: 8080\n"
	got, err := testable.TestLines(fromyaml.FromYaml(), in)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{"server":{"host":"localhost","port":8080}}` {
		t.Fatalf("got %q", got)
	}
}

func TestFromYaml_List(t *testing.T) {
	in := "- a\n- b\n- c\n"
	got, err := testable.TestLines(fromyaml.FromYaml(), in)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `["a","b","c"]` {
		t.Fatalf("got %q", got)
	}
}

func TestFromYaml_EmptyInput(t *testing.T) {
	got, err := testable.TestLines(fromyaml.FromYaml(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `null` {
		t.Fatalf("got %q, want null", got)
	}
}
