package fromtsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// Error is the sentinel error type for the fromtsv package.
type Error string

func (e Error) Error() string { return string(e) }

const (
	// errTSV prefixes a TSV parse failure.
	errTSV Error = "tsv"
	// errJSON prefixes a JSON encoding failure.
	errJSON Error = "json"
)

// marshal is the JSON encoder used to render rows. It is a package variable so
// tests can inject a failing encoder and exercise the error path (the default,
// json.Marshal of a string-keyed map, never fails in production).
var marshal = json.Marshal

// FromTsv returns a command that converts TSV (tab-separated) input into
// newline-delimited JSON: each row becomes one compact JSON object keyed by the
// column headers.
//
// Flags:
//   - FromTSVWithoutHeader: treat every row as data and synthesize col1, col2…
//   - FromTSVTrim: trim leading whitespace in each field
func FromTsv(opts ...any) gloo.Command[[]byte, []byte] {
	f := gloo.NewParameters[gloo.File, flags](opts...).Flags
	noHeader := bool(f.noHeader)
	trim := bool(f.trimSpaces)
	return patterns.Accumulate(func(in [][]byte) ([][]byte, error) {
		return rowsToJSON(in, '\t', noHeader, trim)
	})
}

// parseRecords reads every TSV record from the joined input, tolerating ragged
// rows. It returns nil records (no error) when the input is blank.
func parseRecords(in [][]byte, comma rune, trim bool) ([][]string, error) {
	raw := bytes.Join(in, []byte{'\n'})
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, nil
	}
	r := csv.NewReader(bytes.NewReader(raw))
	r.Comma = comma
	r.TrimLeadingSpace = trim
	r.FieldsPerRecord = -1 // tolerate ragged rows
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errTSV, err)
	}
	return records, nil
}

// syntheticHeaders builds col1, col2, … wide enough for the widest record.
func syntheticHeaders(records [][]string) []string {
	maxCols := 0
	for _, rec := range records {
		if len(rec) > maxCols {
			maxCols = len(rec)
		}
	}
	headers := make([]string, maxCols)
	for i := range headers {
		headers[i] = fmt.Sprintf("col%d", i+1)
	}
	return headers
}

// splitHeader separates header names from data rows. With noHeader, every record
// is data and the headers are synthesized.
func splitHeader(records [][]string, noHeader bool) (headers []string, dataRows [][]string) {
	if noHeader {
		return syntheticHeaders(records), records
	}
	return records[0], records[1:]
}

// recordToJSON encodes one record as a compact JSON object keyed by headers.
// Surplus fields beyond the header count are dropped.
func recordToJSON(headers, rec []string) ([]byte, error) {
	obj := make(map[string]any, len(headers))
	for i, value := range rec {
		if i < len(headers) {
			obj[headers[i]] = value
		}
	}
	enc, err := marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errJSON, err)
	}
	return enc, nil
}

func rowsToJSON(in [][]byte, comma rune, noHeader, trim bool) ([][]byte, error) {
	records, err := parseRecords(in, comma, trim)
	if err != nil || len(records) == 0 {
		return nil, err
	}
	headers, dataRows := splitHeader(records, noHeader)
	out := make([][]byte, 0, len(dataRows))
	for _, rec := range dataRows {
		enc, err := recordToJSON(headers, rec)
		if err != nil {
			return nil, err
		}
		out = append(out, enc)
	}
	return out, nil
}
