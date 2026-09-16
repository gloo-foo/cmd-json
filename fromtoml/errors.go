package fromtoml

import (
	errs "github.com/gomatic/go-error"
)

// The package's error vocabulary: errs.Const sentinels matched with errors.Is.
const (
	// ErrTOML marks a TOML decode failure.
	ErrTOML errs.Const = "toml"
	// ErrJSON marks a JSON encoding failure.
	ErrJSON errs.Const = "json"
)
