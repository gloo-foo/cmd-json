package fromcsv

type (
	Delimiter         rune
	fromCSVNoHeader   bool
	fromCSVTrimSpaces bool
)

const (
	FromCSVWithHeader    fromCSVNoHeader   = false
	FromCSVWithoutHeader fromCSVNoHeader   = true
	FromCSVNoTrim        fromCSVTrimSpaces = false
	FromCSVTrim          fromCSVTrimSpaces = true
)

type flags struct {
	delimiter  rune
	noHeader   fromCSVNoHeader
	trimSpaces fromCSVTrimSpaces
}

func (d Delimiter) Configure(f *flags)         { f.delimiter = rune(d) }
func (n fromCSVNoHeader) Configure(f *flags)   { f.noHeader = n }
func (t fromCSVTrimSpaces) Configure(f *flags) { f.trimSpaces = t }
