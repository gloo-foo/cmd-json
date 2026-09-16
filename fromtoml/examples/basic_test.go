package fromtoml_test

import (
	"fmt"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	"github.com/gloo-foo/cmd-json/fromtoml"
)

func ExampleFromToml() {
	in := "title = \"TOML Example\"\n\n" +
		"[owner]\nname = \"Alice\"\nage = 30\n\n" +
		"[database]\nserver = \"localhost\"\nport = 5432\n"
	lines, _ := testable.TestLines(fromtoml.FromToml(), run.Input(in))
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// {"database":{"port":5432,"server":"localhost"},"owner":{"age":30,"name":"Alice"},"title":"TOML Example"}
}
