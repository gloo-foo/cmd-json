package selectcmd_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	selectcmd "github.com/gloo-foo/cmd-json/select"
)

func ExampleSelect_hasField() {
	in := `{"name":"Alice","age":30}` + "\n" +
		`{"name":"Bob"}` + "\n" +
		`{"name":"Charlie","age":35}` + "\n"
	lines, _ := testable.TestLines(selectcmd.Select(selectcmd.HasField("age")), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"age":30,"name":"Alice"}
	// {"age":35,"name":"Charlie"}
}

func ExampleSelect_fieldEquals() {
	in := `{"name":"Alice","status":"active"}` + "\n" +
		`{"name":"Bob","status":"inactive"}` + "\n" +
		`{"name":"Charlie","status":"active"}` + "\n"
	lines, _ := testable.TestLines(selectcmd.Select(selectcmd.FieldEquals("status", "active")), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"name":"Alice","status":"active"}
	// {"name":"Charlie","status":"active"}
}

func ExampleSelect_and() {
	in := `{"name":"Alice","age":30,"status":"active"}` + "\n" +
		`{"name":"Bob","age":25,"status":"active"}` + "\n" +
		`{"name":"Charlie","age":35,"status":"inactive"}` + "\n"
	cond := selectcmd.And(selectcmd.HasField("age"), selectcmd.FieldEquals("status", "active"))
	lines, _ := testable.TestLines(selectcmd.Select(cond), in)
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"age":30,"name":"Alice","status":"active"}
	// {"age":25,"name":"Bob","status":"active"}
}
