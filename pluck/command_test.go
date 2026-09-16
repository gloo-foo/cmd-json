package pluck_test

import (
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	jsoncmd "github.com/gloo-foo/cmd-json"
	"github.com/gloo-foo/cmd-json/pluck"
)

func TestPluck_KeepsNamedFields(t *testing.T) {
	lines, err := testable.TestLines(pluck.Pluck("name", "age"),
		`{"name":"Alice","age":30,"city":"NYC"}`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"age":30,"name":"Alice"}` {
		t.Fatalf("got %q, want plucked name+age", lines)
	}
}

func TestPluck_MultipleLines(t *testing.T) {
	in := `{"name":"Alice","age":30,"city":"NYC"}` + "\n" +
		`{"name":"Bob","age":25,"city":"LA"}` + "\n"
	lines, err := testable.TestLines(pluck.Pluck("name"), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{`{"name":"Alice"}`, `{"name":"Bob"}`}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

func TestPluck_DropsObjectsWithoutFields(t *testing.T) {
	in := `{"city":"NYC"}` + "\n" + `{"name":"Bob"}` + "\n"
	lines, err := testable.TestLines(pluck.Pluck("name"), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"name":"Bob"}` {
		t.Fatalf("got %q, want only the object with name", lines)
	}
}

func TestPluck_DropsNonObjects(t *testing.T) {
	in := `[1,2,3]` + "\n" + `{"name":"Bob"}` + "\n"
	lines, err := testable.TestLines(pluck.Pluck("name"), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"name":"Bob"}` {
		t.Fatalf("got %q, want only the object", lines)
	}
}

func TestPluck_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(pluck.Pluck("name"), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("got %q, want empty", lines)
	}
}

// TestPluck_AfterDecode is an integration test: a pretty-printed top-level array
// is normalized by Decode, then plucked.
func TestPluck_AfterDecode(t *testing.T) {
	pipe := gloo.Pipe(jsoncmd.Decode(), pluck.Pluck("name", "age"))
	in := "[\n  {\"name\":\"Alice\",\"age\":30,\"city\":\"NYC\"},\n  {\"name\":\"Bob\",\"age\":25}\n]\n"
	lines, err := testable.TestLines(pipe, run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{`{"age":30,"name":"Alice"}`, `{"age":25,"name":"Bob"}`}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}
