package fromyaml

import (
	"bytes"
	"encoding/json"
	"fmt"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
	"gopkg.in/yaml.v3"
)

// Error is the sentinel error type for the fromyaml package.
type Error string

func (e Error) Error() string { return string(e) }

const (
	// errYAML prefixes a YAML decode failure.
	errYAML Error = "yaml"
	// errJSON prefixes a JSON encoding failure (e.g. a YAML mapping with
	// non-string keys, which JSON cannot represent).
	errJSON Error = "json"
)

// yamlToJSON decodes one buffered YAML document and re-encodes it as a single
// compact JSON value.
func yamlToJSON(in [][]byte) ([]byte, error) {
	data := bytes.Join(in, []byte{'\n'})
	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("%w: %w", errYAML, err)
	}
	out, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errJSON, err)
	}
	return out, nil
}

// FromYaml converts a YAML document on the input into a single compact JSON
// value (object, array, or scalar), letting YAML feed the rest of the json
// pipeline. The whole input is one document, so it is buffered and emitted as
// one value.
func FromYaml(opts ...any) gloo.Command[[]byte, []byte] {
	_ = gloo.NewParameters[gloo.File, flags](opts...).Flags
	return patterns.Aggregate(yamlToJSON)
}
