package selectcmd

import (
	gloo "github.com/gloo-foo/framework"

	jsoncmd "github.com/gloo-foo/cmd-json"
)

// Select returns a command that keeps each JSON value for which condition
// reports true, emitting one compact value per line. Values that fail the
// condition are dropped.
//
// Shell analogue: jq 'select(...)'
func Select(condition Condition) gloo.Command[[]byte, []byte] {
	return jsoncmd.Process(func(v jsoncmd.Value) (jsoncmd.Value, bool, error) {
		return v, condition(v), nil
	})
}
