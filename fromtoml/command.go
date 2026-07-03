package fromtoml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// unmarshal is the TOML decoder. It is a package variable for symmetry with
// marshal; tests can swap it if needed.
var unmarshal = toml.Unmarshal

// tomlToJSON decodes one buffered TOML document and re-encodes it as a single
// compact JSON object.
func tomlToJSON(in [][]byte) ([]byte, error) {
	data := bytes.Join(in, []byte{'\n'})
	result := map[string]any{}
	if err := unmarshal(data, &result); err != nil {
		return nil, ErrTOML.With(err)
	}
	out, err := marshal(result)
	if err != nil {
		return nil, ErrJSON.With(err)
	}
	return out, nil
}

// FromToml converts a TOML document on the input into a single compact JSON
// object, letting TOML feed the rest of the json pipeline. The whole input is
// one document, so it is buffered and emitted as one value.
func FromToml(opts ...any) gloo.Command[[]byte, []byte] {
	_ = gloo.NewParameters[gloo.File, flags](opts...).Flags
	return patterns.Aggregate(tomlToJSON)
}
