package fromtsv

// Error is the sentinel error type for the fromtsv package.
type Error string

func (e Error) Error() string { return string(e) }

const (
	// errTSV prefixes a TSV parse failure.
	errTSV Error = "tsv"
	// errJSON prefixes a JSON encoding failure.
	errJSON Error = "json"
)
