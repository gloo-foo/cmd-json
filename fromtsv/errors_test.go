package fromtsv

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable"
)

func TestErrTSV_Text(t *testing.T) {
	if ErrTSV.Error() != "tsv" {
		t.Fatalf("got %q, want %q", ErrTSV.Error(), "tsv")
	}
}

func TestFromTsv_MalformedTSVErrors(t *testing.T) {
	// An unterminated quoted field is a parse error.
	_, err := testable.TestLines(FromTsv(), "name\n\"unterminated\n")
	if !errors.Is(err, ErrTSV) {
		t.Fatalf("got %v, want ErrTSV", err)
	}
}
