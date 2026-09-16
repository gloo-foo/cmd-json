package fromtsv

import (
	errs "github.com/gomatic/go-error"
)

// The package's error vocabulary: errs.Const sentinels matched with errors.Is.
const (
	// ErrTSV marks a TSV parse failure.
	ErrTSV errs.Const = "tsv"
	// ErrJSON marks a JSON encoding failure.
	ErrJSON errs.Const = "json"
)
