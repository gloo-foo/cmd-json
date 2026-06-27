package pluck_test

import (
	"fmt"

	"github.com/gloo-foo/cmd-json/pluck"
	"github.com/gloo-foo/testable"
)

func ExamplePluck() {
	// echo '{...}\n{...}' | json pluck name age
	in := `{"name":"Alice","age":30,"city":"NYC"}` + "\n" +
		`{"name":"Bob","age":25,"city":"LA"}` + "\n"
	lines, _ := testable.TestLines(pluck.Pluck("name", "age"), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"age":30,"name":"Alice"}
	// {"age":25,"name":"Bob"}
}

func ExamplePluck_singleField() {
	in := `{"name":"Alice","age":30,"city":"NYC"}` + "\n" +
		`{"name":"Bob","age":25,"city":"LA"}` + "\n"
	lines, _ := testable.TestLines(pluck.Pluck("name"), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"name":"Alice"}
	// {"name":"Bob"}
}
