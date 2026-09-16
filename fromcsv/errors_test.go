package fromcsv

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

func TestErrCSV_Text(t *testing.T) {
	if ErrCSV.Error() != "csv" {
		t.Fatalf("got %q, want %q", ErrCSV.Error(), "csv")
	}
}

func TestFromCsv_MalformedCSVErrors(t *testing.T) {
	// An unterminated quoted field is a CSV parse error.
	_, err := testable.TestLines(FromCsv(), "name\n\"unterminated\n")
	if !errors.Is(err, ErrCSV) {
		t.Fatalf("got %v, want ErrCSV", err)
	}
}
