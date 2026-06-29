package fromtsv_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	"github.com/gloo-foo/cmd-json/fromtsv"
)

func ExampleFromTsv() {
	in := "name\tage\tcity\nAlice\t30\tNYC\nBob\t25\tLA\n"
	lines, _ := testable.TestLines(fromtsv.FromTsv(), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"age":"30","city":"NYC","name":"Alice"}
	// {"age":"25","city":"LA","name":"Bob"}
}

func ExampleFromTsv_noHeader() {
	in := "Alice\t30\tNYC\nBob\t25\tLA\n"
	lines, _ := testable.TestLines(fromtsv.FromTsv(fromtsv.FromTSVWithoutHeader), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"col1":"Alice","col2":"30","col3":"NYC"}
	// {"col1":"Bob","col2":"25","col3":"LA"}
}
