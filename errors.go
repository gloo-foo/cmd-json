package command

// Error is the sentinel error type for the json command package. Every error
// the package emits is a constant of this type, so callers can match it with
// errors.Is rather than comparing strings.
type Error string

func (e Error) Error() string { return string(e) }
