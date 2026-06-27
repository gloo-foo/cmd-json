package fromcsv

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

func TestError_Error(t *testing.T) {
	if errCSV.Error() != "csv" {
		t.Fatalf("got %q, want %q", errCSV.Error(), "csv")
	}
}

func TestFromCsv_Trim(t *testing.T) {
	// Leading whitespace in each field is removed when FromCSVTrim is set.
	got, err := testable.TestLines(FromCsv(FromCSVTrim), "name, age\nAlice,  30\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != `{"age":"30","name":"Alice"}` {
		t.Fatalf("got %q, want trimmed value", got)
	}
}

func TestFromCsv_MalformedCSVErrors(t *testing.T) {
	// An unterminated quoted field is a CSV parse error.
	_, err := testable.TestLines(FromCsv(), "name\n\"unterminated\n")
	if !errors.Is(err, errCSV) {
		t.Fatalf("got %v, want errCSV", err)
	}
}

func TestFromCsv_MarshalError(t *testing.T) {
	original := marshal
	marshal = func(any) ([]byte, error) { return nil, errJSON }
	t.Cleanup(func() { marshal = original })

	_, err := testable.TestLines(FromCsv(), "name\nAlice\n")
	if !errors.Is(err, errJSON) {
		t.Fatalf("got %v, want errJSON", err)
	}
}
