package selectcmd_test

import (
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"

	jsoncmd "github.com/gloo-foo/cmd-json"
	selectcmd "github.com/gloo-foo/cmd-json/select"
)

func lines(t *testing.T, cmd gloo.Command[[]byte, []byte], in string) []string {
	t.Helper()
	got, err := testable.TestLines(cmd, in)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func eq(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

const sample = `{"name":"Alice","age":30,"status":"active"}` + "\n" +
	`{"name":"Bob","age":25,"status":"inactive"}` + "\n" +
	`{"name":"Charlie","status":"active"}` + "\n"

func TestSelect_HasField(t *testing.T) {
	got := lines(t, selectcmd.Select(selectcmd.HasField("age")), sample)
	eq(t, got, []string{
		`{"age":30,"name":"Alice","status":"active"}`,
		`{"age":25,"name":"Bob","status":"inactive"}`,
	})
}

func TestSelect_FieldEquals(t *testing.T) {
	got := lines(t, selectcmd.Select(selectcmd.FieldEquals("status", "active")), sample)
	eq(t, got, []string{
		`{"age":30,"name":"Alice","status":"active"}`,
		`{"name":"Charlie","status":"active"}`,
	})
}

func TestSelect_And(t *testing.T) {
	cond := selectcmd.And(selectcmd.HasField("age"), selectcmd.FieldEquals("status", "active"))
	got := lines(t, selectcmd.Select(cond), sample)
	eq(t, got, []string{`{"age":30,"name":"Alice","status":"active"}`})
}

func TestSelect_Or(t *testing.T) {
	cond := selectcmd.Or(selectcmd.FieldEquals("name", "Bob"), selectcmd.FieldEquals("name", "Charlie"))
	got := lines(t, selectcmd.Select(cond), sample)
	eq(t, got, []string{
		`{"age":25,"name":"Bob","status":"inactive"}`,
		`{"name":"Charlie","status":"active"}`,
	})
}

func TestSelect_Not(t *testing.T) {
	got := lines(t, selectcmd.Select(selectcmd.Not(selectcmd.HasField("age"))), sample)
	eq(t, got, []string{`{"name":"Charlie","status":"active"}`})
}

func TestSelect_FieldMatches(t *testing.T) {
	cond := selectcmd.FieldMatches("age", func(v any) bool {
		age, ok := v.(float64)
		return ok && age > 28
	})
	got := lines(t, selectcmd.Select(cond), sample)
	eq(t, got, []string{`{"age":30,"name":"Alice","status":"active"}`})
}

func TestSelect_EmptyInput(t *testing.T) {
	got := lines(t, selectcmd.Select(selectcmd.HasField("age")), "")
	if len(got) != 0 {
		t.Fatalf("got %q, want empty", got)
	}
}

// TestSelect_AfterDecode is an integration test composing Decode with Select.
func TestSelect_AfterDecode(t *testing.T) {
	pipe := gloo.Pipe(jsoncmd.Decode(), selectcmd.Select(selectcmd.HasField("age")))
	in := "[\n  {\"name\":\"Alice\",\"age\":30},\n  {\"name\":\"Charlie\"}\n]\n"
	got := lines(t, pipe, in)
	eq(t, got, []string{`{"age":30,"name":"Alice"}`})
}
