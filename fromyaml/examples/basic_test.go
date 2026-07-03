package fromyaml_test

import (
	"fmt"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	"github.com/gloo-foo/cmd-json/fromyaml"
)

func ExampleFromYaml() {
	in := "title: YAML Example\nowner:\n  name: Alice\n  age: 30\n"
	lines, _ := testable.TestLines(fromyaml.FromYaml(), run.Input(in))
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"owner":{"age":30,"name":"Alice"},"title":"YAML Example"}
}
