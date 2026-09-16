package fromcsv

import (
	errs "github.com/gomatic/go-error"
)

// The package's error vocabulary: errs.Const sentinels matched with errors.Is.
const (
	// ErrCSV marks a CSV parse failure.
	ErrCSV errs.Const = "csv"
	// ErrJSON marks a JSON encoding failure.
	ErrJSON errs.Const = "json"
)
