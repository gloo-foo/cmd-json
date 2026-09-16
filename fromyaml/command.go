package fromyaml

import (
	"bytes"
	"encoding/json"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
	"gopkg.in/yaml.v3"
)

// yamlToJSON decodes one buffered YAML document and re-encodes it as a single
// compact JSON value.
func yamlToJSON(in [][]byte) ([]byte, error) {
	data := bytes.Join(in, []byte{'\n'})
	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, ErrYAML.With(err)
	}
	out, err := json.Marshal(result)
	if err != nil {
		return nil, ErrJSON.With(err)
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
