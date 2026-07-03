package command

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
	errs "github.com/gomatic/go-error"
)

// errMarshalFail is returned by the injected encoder to drive every marshal
// error path.
const errMarshalFail errs.Const = "marshal failed"

// withFailingMarshal swaps the package marshal function for one that always
// fails, restoring the original when the test ends.
func withFailingMarshal(t *testing.T) {
	t.Helper()
	original := marshal
	marshal = func(any) ([]byte, error) { return nil, errMarshalFail }
	t.Cleanup(func() { marshal = original })
}

// The marshal error branches are unreachable with the default encoder: values
// re-emitted by JSON, Process, and Decode all originate from a successful
// json.Unmarshal and so always marshal cleanly. Injecting a failing encoder is
// the only honest way to exercise these branches.

func TestJson_MarshalError(t *testing.T) {
	withFailingMarshal(t)
	_, err := testable.TestLines(JSON(), `{"a":1}`+"\n")
	if !errors.Is(err, ErrMarshal) {
		t.Fatalf("got %v, want ErrMarshal", err)
	}
	if !errors.Is(err, errMarshalFail) {
		t.Fatalf("got %v, want the injected cause wrapped", err)
	}
}

func TestProcess_MarshalError(t *testing.T) {
	withFailingMarshal(t)
	keep := func(v Value) (Value, bool, error) { return v, true, nil }
	_, err := testable.TestLines(Process(keep), `{"a":1}`+"\n")
	if !errors.Is(err, ErrMarshal) {
		t.Fatalf("got %v, want ErrMarshal", err)
	}
}

func TestDecode_MarshalError(t *testing.T) {
	withFailingMarshal(t)
	_, err := testable.TestLines(Decode(), `{"a":1}`+"\n")
	if !errors.Is(err, ErrMarshal) {
		t.Fatalf("got %v, want ErrMarshal", err)
	}
}
