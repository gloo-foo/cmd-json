package command

import (
	"encoding/json"
	"fmt"
)

// errMarshal prefixes a re-encoding failure.
const errMarshal Error = "json: marshal"

// marshal is the JSON encoder used to re-emit values. It is a package variable
// so tests can inject a failing encoder and exercise the error path. Values
// re-emitted here always originate from a successful json.Unmarshal, so the
// default json.Marshal never fails in production.
var marshal = json.Marshal

// encodeValue marshals v to compact JSON, wrapping any failure in errMarshal.
func encodeValue(v Value) ([]byte, error) {
	enc, err := marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errMarshal, err)
	}
	return enc, nil
}
