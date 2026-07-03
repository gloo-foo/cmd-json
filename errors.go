package command

import (
	errs "github.com/gomatic/go-error"
)

// The package's error vocabulary. Every error the package emits is one of
// these errs.Const sentinels (wrapped around its cause with Const.With), so
// callers match failure classes with errors.Is rather than comparing strings.
const (
	// ErrInvalidJSON marks a parse failure on a single JSON input line.
	ErrInvalidJSON errs.Const = "json: invalid JSON"
	// ErrInvalidInput marks a decode failure on a single Process line.
	ErrInvalidInput errs.Const = "json: invalid input"
	// ErrDecode marks a streaming decode failure in Decode.
	ErrDecode errs.Const = "json: decode"
	// ErrMarshal marks a re-encoding failure.
	ErrMarshal errs.Const = "json: marshal"
	// ErrQuery marks a cirql query failure (parse or run).
	ErrQuery errs.Const = "json: query"
)
