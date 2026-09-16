package command_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	command "github.com/gloo-foo/cmd-json"
)

func TestJson_Passthrough(t *testing.T) {
	lines, err := testable.TestLines(command.JSON(), `{"name":"alice","age":30}`+"\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	// json.Marshal sorts keys, so output is deterministic
	want := `{"age":30,"name":"alice"}`
	if lines[0] != want {
		t.Fatalf("got %q, want %q", lines[0], want)
	}
}

func TestJson_MultipleLines(t *testing.T) {
	input := `{"a":1}` + "\n" + `{"b":2}` + "\n"
	lines, err := testable.TestLines(command.JSON(), run.Input(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	if lines[0] != `{"a":1}` {
		t.Errorf("line 0: got %q, want %q", lines[0], `{"a":1}`)
	}
	if lines[1] != `{"b":2}` {
		t.Errorf("line 1: got %q, want %q", lines[1], `{"b":2}`)
	}
}

func TestJson_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(command.JSON(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("got %q, want empty", lines)
	}
}

func TestJson_CompactsWhitespace(t *testing.T) {
	input := `{ "key" : "value" , "num" : 42 }` + "\n"
	lines, err := testable.TestLines(command.JSON(), run.Input(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	want := `{"key":"value","num":42}`
	if lines[0] != want {
		t.Fatalf("got %q, want %q", lines[0], want)
	}
}

func TestJson_Array(t *testing.T) {
	input := `[1,2,3]` + "\n"
	lines, err := testable.TestLines(command.JSON(), run.Input(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `[1,2,3]` {
		t.Fatalf("got %q, want [[1,2,3]]", lines)
	}
}

func TestJson_Scalar(t *testing.T) {
	input := `"hello"` + "\n"
	lines, err := testable.TestLines(command.JSON(), run.Input(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != `"hello"` {
		t.Fatalf("got %q, want [\"hello\"]", lines)
	}
}

func TestJson_InvalidJSON(t *testing.T) {
	_, err := testable.TestLines(command.JSON(), "not json\n")
	if !errors.Is(err, command.ErrInvalidJSON) {
		t.Fatalf("got %v, want ErrInvalidJSON", err)
	}
}

func ExampleJSON() {
	lines, _ := testable.TestLines(command.JSON(), `{"b":2,"a":1}`+"\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"a":1,"b":2}
}
