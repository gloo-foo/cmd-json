package fromtsv_test

import (
	"testing"

	"github.com/gloo-foo/cmd-json/fromtsv"
	"github.com/gloo-foo/testable"
)

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

func TestFromTsv_Basic(t *testing.T) {
	in := "name\tage\tcity\nAlice\t30\tNYC\nBob\t25\tLA\n"
	got, err := testable.TestLines(fromtsv.FromTsv(), in)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, got, []string{
		`{"age":"30","city":"NYC","name":"Alice"}`,
		`{"age":"25","city":"LA","name":"Bob"}`,
	})
}

func TestFromTsv_NoHeader(t *testing.T) {
	in := "Alice\t30\tNYC\nBob\t25\tLA\n"
	got, err := testable.TestLines(fromtsv.FromTsv(fromtsv.FromTSVWithoutHeader), in)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, got, []string{
		`{"col1":"Alice","col2":"30","col3":"NYC"}`,
		`{"col1":"Bob","col2":"25","col3":"LA"}`,
	})
}

func TestFromTsv_EmptyInput(t *testing.T) {
	got, err := testable.TestLines(fromtsv.FromTsv(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("got %q, want empty", got)
	}
}
