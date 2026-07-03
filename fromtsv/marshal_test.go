package fromtsv

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

func TestFromTsv_MarshalError(t *testing.T) {
	original := marshal
	marshal = func(any) ([]byte, error) { return nil, ErrJSON }
	t.Cleanup(func() { marshal = original })

	_, err := testable.TestLines(FromTsv(), "name\nAlice\n")
	if !errors.Is(err, ErrJSON) {
		t.Fatalf("got %v, want ErrJSON", err)
	}
}
