package fromtsv

type (
	// fromTSVNoHeader selects whether the first record is data (headerless) or
	// the column headers.
	fromTSVNoHeader bool
	// fromTSVTrimSpaces selects whether leading whitespace is trimmed from each
	// field.
	fromTSVTrimSpaces bool
)

const (
	FromTSVWithHeader    fromTSVNoHeader   = false
	FromTSVWithoutHeader fromTSVNoHeader   = true
	FromTSVNoTrim        fromTSVTrimSpaces = false
	FromTSVTrim          fromTSVTrimSpaces = true
)

// flags is the option set folded from a FromTsv call's option values.
type flags struct {
	isHeaderless fromTSVNoHeader
	shouldTrim   fromTSVTrimSpaces
}

// with folds one option value into the flag set. Values of any other type are
// ignored: FromTsv reads only the upstream stream.
func (f flags) with(o any) flags {
	switch v := o.(type) {
	case fromTSVNoHeader:
		f.isHeaderless = v
	case fromTSVTrimSpaces:
		f.shouldTrim = v
	}
	return f
}

// fold collapses the FromTsv option values into the flag set.
func fold(opts []any) flags {
	var f flags
	for _, o := range opts {
		f = f.with(o)
	}
	return f
}
