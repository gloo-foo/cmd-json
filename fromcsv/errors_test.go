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

func TestFromCsv_MalformedCSVErrors(t *testing.T) {
	// An unterminated quoted field is a CSV parse error.
	_, err := testable.TestLines(FromCsv(), "name\n\"unterminated\n")
	if !errors.Is(err, errCSV) {
		t.Fatalf("got %v, want errCSV", err)
	}
}
