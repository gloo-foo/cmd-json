package command_test

import (
	"fmt"
	"strings"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	command "github.com/gloo-foo/cmd-json"
)

// errBoom is a stand-in error a test Processor returns to exercise the
// error-propagation path of Process.
const errBoom command.Error = "boom"

// keepAdults keeps objects whose "age" is >= 18, unchanged.
func keepAdults(v command.Value) (command.Value, bool, error) {
	obj, ok := command.AsMap(v)
	if !ok {
		return nil, false, nil
	}
	age, ok := obj["age"].(float64)
	return v, ok && age >= 18, nil
}

// renameToUpper rewrites the "name" field to upper case.
func renameToUpper(v command.Value) (command.Value, bool, error) {
	obj, ok := command.AsMap(v)
	if !ok {
		return v, true, nil
	}
	if name, ok := obj["name"].(string); ok {
		obj["name"] = strings.ToUpper(name)
	}
	return obj, true, nil
}

func TestProcess_Filters(t *testing.T) {
	in := `{"name":"al","age":30}` + "\n" + `{"name":"bo","age":12}` + "\n"
	lines, err := testable.TestLines(command.Process(keepAdults), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"age":30,"name":"al"}` {
		t.Fatalf("got %q, want one adult object", lines)
	}
}

func TestProcess_Transforms(t *testing.T) {
	lines, err := testable.TestLines(command.Process(renameToUpper), `{"name":"alice"}`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"name":"ALICE"}` {
		t.Fatalf("got %q, want renamed object", lines)
	}
}

func TestProcess_SkipsBlankLines(t *testing.T) {
	in := `{"age":40}` + "\n\n   \n" + `{"age":50}` + "\n"
	lines, err := testable.TestLines(command.Process(keepAdults), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), lines)
	}
}

func TestProcess_InvalidJSONErrors(t *testing.T) {
	if _, err := testable.TestLines(command.Process(keepAdults), "not json\n"); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestProcess_PropagatesProcessorError(t *testing.T) {
	boom := func(command.Value) (command.Value, bool, error) {
		return nil, false, errBoom
	}
	if _, err := testable.TestLines(command.Process(boom), `{"a":1}`+"\n"); err == nil {
		t.Fatal("expected the Processor error to propagate")
	}
}

func TestDecode_InvalidJSONErrors(t *testing.T) {
	if _, err := testable.TestLines(command.Decode(), "{not valid\n"); err == nil {
		t.Fatal("expected error for malformed JSON document")
	}
}

func TestProcess_DropsNonObjects(t *testing.T) {
	in := `42` + "\n" + `{"age":21}` + "\n"
	lines, err := testable.TestLines(command.Process(keepAdults), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"age":21}` {
		t.Fatalf("got %q, want only the object", lines)
	}
}

func TestDecode_JSONL(t *testing.T) {
	in := `{"a":1}` + "\n" + `{"b":2}` + "\n"
	lines, err := testable.TestLines(command.Decode(), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || lines[0] != `{"a":1}` || lines[1] != `{"b":2}` {
		t.Fatalf("got %q", lines)
	}
}

func TestDecode_StreamsTopLevelArray(t *testing.T) {
	lines, err := testable.TestLines(command.Decode(), `[{"a":1},{"b":2},3]`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{`{"a":1}`, `{"b":2}`, `3`}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

func TestDecode_PrettyDocument(t *testing.T) {
	in := "{\n  \"name\": \"alice\",\n  \"age\": 30\n}\n"
	lines, err := testable.TestLines(command.Decode(), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"age":30,"name":"alice"}` {
		t.Fatalf("got %q, want one compact object", lines)
	}
}

func TestDecode_Concatenated(t *testing.T) {
	lines, err := testable.TestLines(command.Decode(), `{"a":1} {"b":2}`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), lines)
	}
}

func TestDecode_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(command.Decode(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("got %q, want empty", lines)
	}
}

func TestAsMap(t *testing.T) {
	if m, ok := command.AsMap(map[string]any{"x": 1.0}); !ok || m["x"] != 1.0 {
		t.Fatalf("AsMap on object failed: %v %v", m, ok)
	}
	if _, ok := command.AsMap("scalar"); ok {
		t.Fatal("AsMap on non-object should be false")
	}
}

// TestIntegration_DecodeThenProcess wires Decode into Process: a pretty-printed
// top-level array is normalized to JSONL, then filtered.
func TestIntegration_DecodeThenProcess(t *testing.T) {
	pipe := gloo.Pipe(command.Decode(), command.Process(keepAdults))
	in := "[\n  {\"name\":\"al\",\"age\":30},\n  {\"name\":\"bo\",\"age\":10}\n]\n"
	lines, err := testable.TestLines(pipe, run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `{"age":30,"name":"al"}` {
		t.Fatalf("got %q, want one adult", lines)
	}
}

func ExampleProcess() {
	upper := command.Process(func(v command.Value) (command.Value, bool, error) {
		obj, _ := command.AsMap(v)
		return obj, true, nil
	})
	lines, _ := testable.TestLines(upper, `{"b":2,"a":1}`+"\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"a":1,"b":2}
}

func ExampleDecode() {
	lines, _ := testable.TestLines(command.Decode(), `[{"a":1},{"b":2}]`+"\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"a":1}
	// {"b":2}
}
