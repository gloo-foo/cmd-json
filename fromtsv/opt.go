package fromtsv

type (
	fromTSVNoHeader   bool
	fromTSVTrimSpaces bool
)

const (
	FromTSVWithHeader    fromTSVNoHeader   = false
	FromTSVWithoutHeader fromTSVNoHeader   = true
	FromTSVNoTrim        fromTSVTrimSpaces = false
	FromTSVTrim          fromTSVTrimSpaces = true
)

type flags struct {
	noHeader   fromTSVNoHeader
	trimSpaces fromTSVTrimSpaces
}

func (n fromTSVNoHeader) Configure(f *flags)   { f.noHeader = n }
func (t fromTSVTrimSpaces) Configure(f *flags) { f.trimSpaces = t }
