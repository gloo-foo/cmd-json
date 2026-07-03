package command_test

import (
	"testing"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	command "github.com/gloo-foo/cmd-json"
)

func TestQuery_FilterMapSortLimit(t *testing.T) {
	in := `{"name":"a","stars":500}` + "\n" +
		`{"name":"b","stars":3000}` + "\n" +
		`{"name":"c","stars":2000}` + "\n"
	q := command.QueryScript(`filter .stars > 1000 | map { name: .name, stars: .stars } | sort .stars desc | limit 2`)
	lines, err := testable.TestLines(command.Query(q), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{`{"name":"b","stars":3000}`, `{"name":"c","stars":2000}`}
	if len(lines) != 2 || lines[0] != want[0] || lines[1] != want[1] {
		t.Fatalf("got %q, want %q", lines, want)
	}
}

func TestQuery_TopLevelArrayInput(t *testing.T) {
	in := `[{"n":1},{"n":2},{"n":3}]`
	lines, err := testable.TestLines(command.Query(`filter .n >= 2`), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("got %q, want 2 rows", lines)
	}
}

func TestQuery_Reduce(t *testing.T) {
	in := `{"v":2}` + "\n" + `{"v":4}` + "\n"
	lines, err := testable.TestLines(command.Query(`reduce sum(.v)`), run.Input(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "6" {
		t.Fatalf("got %q, want [6]", lines)
	}
}

func TestQuery_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(command.Query(`filter .a`), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("got %q, want no output", lines)
	}
}

func TestQuery_ParseErrorFailsStream(t *testing.T) {
	if _, err := testable.TestLines(command.Query(`map {`), `{"a":1}`+"\n"); err == nil {
		t.Fatal("expected an invalid query to fail the stream")
	}
}

func TestQuery_DecodeErrorFailsStream(t *testing.T) {
	if _, err := testable.TestLines(command.Query(`filter .a`), "not json\n"); err == nil {
		t.Fatal("expected decode error to propagate")
	}
}

func TestQuery_RunErrorFailsStream(t *testing.T) {
	in := `{"n":"x"}` + "\n" + `{"n":1}` + "\n"
	if _, err := testable.TestLines(command.Query(`sort .n`), run.Input(in)); err == nil {
		t.Fatal("expected an incomparable sort to fail the stream")
	}
}
