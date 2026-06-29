package fromcsv_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	"github.com/gloo-foo/cmd-json/fromcsv"
)

func ExampleFromCsv() {
	in := "name,age,city\nAlice,30,NYC\nBob,25,LA\nCharlie,35,London\n"
	lines, _ := testable.TestLines(fromcsv.FromCsv(), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"age":"30","city":"NYC","name":"Alice"}
	// {"age":"25","city":"LA","name":"Bob"}
	// {"age":"35","city":"London","name":"Charlie"}
}

func ExampleFromCsv_noHeader() {
	in := "Alice,30,NYC\nBob,25,LA\n"
	lines, _ := testable.TestLines(fromcsv.FromCsv(fromcsv.FromCSVWithoutHeader), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"col1":"Alice","col2":"30","col3":"NYC"}
	// {"col1":"Bob","col2":"25","col3":"LA"}
}

func ExampleFromCsv_customDelimiter() {
	in := "name|age|city\nAlice|30|NYC\nBob|25|LA\n"
	lines, _ := testable.TestLines(fromcsv.FromCsv(fromcsv.Delimiter('|')), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"age":"30","city":"NYC","name":"Alice"}
	// {"age":"25","city":"LA","name":"Bob"}
}
