package fromcsv

// Error is the sentinel error type for the fromcsv package.
type Error string

func (e Error) Error() string { return string(e) }

const (
	// errCSV prefixes a CSV parse failure.
	errCSV Error = "csv"
	// errJSON prefixes a JSON encoding failure.
	errJSON Error = "json"
)
