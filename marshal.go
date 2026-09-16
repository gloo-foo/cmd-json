package command

import (
	"encoding/json"
)

// marshal is the JSON encoder used to re-emit values. It is a package variable
// so tests can inject a failing encoder and exercise the error path. Values
// re-emitted here always originate from a successful json.Unmarshal, so the
// default json.Marshal never fails in production.
var marshal = json.Marshal

// encodeValue marshals v to compact JSON, wrapping any failure in ErrMarshal.
func encodeValue(v Value) ([]byte, error) {
	enc, err := marshal(v)
	if err != nil {
		return nil, ErrMarshal.With(err)
	}
	return enc, nil
}
