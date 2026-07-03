package fromcsv

type (
	// Delimiter is the CSV field delimiter (default ',').
	Delimiter rune
	// fromCSVNoHeader selects whether the first record is data (headerless) or
	// the column headers.
	fromCSVNoHeader bool
	// fromCSVTrimSpaces selects whether leading whitespace is trimmed from each
	// field.
	fromCSVTrimSpaces bool
)

const (
	FromCSVWithHeader    fromCSVNoHeader   = false
	FromCSVWithoutHeader fromCSVNoHeader   = true
	FromCSVNoTrim        fromCSVTrimSpaces = false
	FromCSVTrim          fromCSVTrimSpaces = true
)

// flags is the option set folded from a FromCsv call's option values.
type flags struct {
	delimiter    Delimiter
	isHeaderless fromCSVNoHeader
	shouldTrim   fromCSVTrimSpaces
}

// with folds one option value into the flag set. Values of any other type are
// ignored: FromCsv reads only the upstream stream.
func (f flags) with(o any) flags {
	switch v := o.(type) {
	case Delimiter:
		f.delimiter = v
	case fromCSVNoHeader:
		f.isHeaderless = v
	case fromCSVTrimSpaces:
		f.shouldTrim = v
	}
	return f
}

// fold collapses the FromCsv option values into the flag set.
func fold(opts []any) flags {
	var f flags
	for _, o := range opts {
		f = f.with(o)
	}
	return f
}
