package command_test

import (
	"errors"
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

func TestQuery_FailurePaths(t *testing.T) {
	tests := []struct {
		wantErr error
		name    string
		query   command.QueryScript
		input   string
	}{
		{
			name:    "parse error fails the stream",
			query:   `map {`,
			input:   `{"a":1}` + "\n",
			wantErr: command.ErrQuery,
		},
		{
			name:    "decode error propagates",
			query:   `filter .a`,
			input:   "not json\n",
			wantErr: command.ErrDecode,
		},
		{
			name:    "run error (incomparable sort) fails the stream",
			query:   `sort .n`,
			input:   `{"n":"x"}` + "\n" + `{"n":1}` + "\n",
			wantErr: command.ErrQuery,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testable.TestLines(command.Query(tt.query), run.Input(tt.input))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
