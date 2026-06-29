package command

import (
	"encoding/json"
	"fmt"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// errInvalidJSON prefixes a parse failure on a single JSON input line.
const errInvalidJSON Error = "json: invalid JSON"

// compact parses one line as JSON and re-emits it in compact, key-sorted form.
func compact(line []byte) ([]byte, error) {
	var v any
	if err := json.Unmarshal(line, &v); err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidJSON, err)
	}
	return encodeValue(v)
}

// JSON returns a command that parses each input line as JSON and re-emits it
// in compact form. Each input line must be valid JSON.
func JSON(opts ...any) gloo.Command[[]byte, []byte] {
	_ = gloo.NewParameters[gloo.File, flags](opts...).Flags
	return patterns.Map(compact)
}
