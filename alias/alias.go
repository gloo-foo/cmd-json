// Package alias provides unprefixed names for the json command's public API.
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-json"
)

// JSON re-exports the command.JSON constructor.
func JSON(opts ...any) gloo.Command[[]byte, []byte] { return command.JSON(opts...) }

// Query re-exports the command.Query constructor — run a cirql pipeline query.
func Query(q command.QueryScript) gloo.Command[[]byte, []byte] { return command.Query(q) }

// QueryScript is the source text of a cirql pipeline query.
type QueryScript = command.QueryScript
