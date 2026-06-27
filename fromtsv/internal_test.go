package fromtsv

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

func TestError_Error(t *testing.T) {
	if errTSV.Error() != "tsv" {
		t.Fatalf("got %q, want %q", errTSV.Error(), "tsv")
	}
}

func TestFromTsv_Trim(t *testing.T) {
	// Leading whitespace in each field is removed when FromTSVTrim is set.
	got, err := testable.TestLines(FromTsv(FromTSVTrim), "name\tage\nAlice\t  30\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{"age":"30","name":"Alice"}` {
		t.Fatalf("got %q, want trimmed value", got)
	}
}

func TestFromTsv_MalformedTSVErrors(t *testing.T) {
	// An unterminated quoted field is a parse error.
	_, err := testable.TestLines(FromTsv(), "name\n\"unterminated\n")
	if !errors.Is(err, errTSV) {
		t.Fatalf("got %v, want errTSV", err)
	}
}

func TestFromTsv_MarshalError(t *testing.T) {
	original := marshal
	marshal = func(any) ([]byte, error) { return nil, errJSON }
	t.Cleanup(func() { marshal = original })

	_, err := testable.TestLines(FromTsv(), "name\nAlice\n")
	if !errors.Is(err, errJSON) {
		t.Fatalf("got %v, want errJSON", err)
	}
}
