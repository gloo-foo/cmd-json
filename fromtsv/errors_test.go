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

func TestFromTsv_MalformedTSVErrors(t *testing.T) {
	// An unterminated quoted field is a parse error.
	_, err := testable.TestLines(FromTsv(), "name\n\"unterminated\n")
	if !errors.Is(err, errTSV) {
		t.Fatalf("got %v, want errTSV", err)
	}
}
