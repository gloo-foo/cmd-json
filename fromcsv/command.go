package fromcsv

import (
	"bytes"
	"encoding/csv"
	"fmt"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// defaultDelimiter is the CSV field delimiter used when Delimiter is not
// supplied.
const defaultDelimiter Delimiter = ','

// FromCsv returns a command that converts CSV input into newline-delimited JSON:
// each row becomes one compact JSON object keyed by the column headers.
//
// Flags:
//   - Delimiter(r): field delimiter (default ',')
//   - FromCSVWithoutHeader: treat every row as data and synthesize col1, col2…
//   - FromCSVTrim: trim leading whitespace in each field
//
// Input is buffered so the CSV reader can honor quoted fields that span lines.
func FromCsv(opts ...any) gloo.Command[[]byte, []byte] {
	f := fold(opts)
	if f.delimiter == 0 {
		f.delimiter = defaultDelimiter
	}
	return patterns.Accumulate(func(in [][]byte) ([][]byte, error) {
		return rowsToJSON(in, f)
	})
}

// parseRecords reads every CSV record from the joined input, tolerating ragged
// rows. It returns nil records (no error) when the input is blank.
func parseRecords(in [][]byte, f flags) ([][]string, error) {
	raw := bytes.Join(in, []byte{'\n'})
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, nil
	}
	r := csv.NewReader(bytes.NewReader(raw))
	r.Comma = rune(f.delimiter)
	r.TrimLeadingSpace = bool(f.shouldTrim)
	r.FieldsPerRecord = -1 // tolerate ragged rows
	records, err := r.ReadAll()
	if err != nil {
		return nil, ErrCSV.With(err)
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

// splitHeader separates header names from data rows. When isHeaderless, every
// record is data and the headers are synthesized.
func splitHeader(records [][]string, isHeaderless fromCSVNoHeader) (headers []string, dataRows [][]string) {
	if bool(isHeaderless) {
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
		return nil, ErrJSON.With(err)
	}
	return enc, nil
}

// rowsToJSON parses the buffered input per the folded flags and encodes each
// data row as its own compact-JSON line.
func rowsToJSON(in [][]byte, f flags) ([][]byte, error) {
	records, err := parseRecords(in, f)
	if err != nil || len(records) == 0 {
		return nil, err
	}
	headers, dataRows := splitHeader(records, f.isHeaderless)
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
