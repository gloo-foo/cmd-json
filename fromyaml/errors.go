package fromyaml

import (
	errs "github.com/gomatic/go-error"
)

// The package's error vocabulary: errs.Const sentinels matched with errors.Is.
const (
	// ErrYAML marks a YAML decode failure.
	ErrYAML errs.Const = "yaml"
	// ErrJSON marks a JSON encoding failure (e.g. a YAML mapping with
	// non-string keys, which JSON cannot represent).
	ErrJSON errs.Const = "json"
)
