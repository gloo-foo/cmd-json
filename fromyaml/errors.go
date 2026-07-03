package fromyaml

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
