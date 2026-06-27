package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// Value is a decoded JSON value: an object (map[string]Value), array ([]Value),
// string, float64, bool, or nil. It is the unit every json command operates on.
type Value = any

// Processor transforms a single decoded JSON value. It returns the (possibly
// rewritten) value and whether to keep it: returning keep=false drops the value
// from the output stream (filter semantics). Returning a non-nil error stops the
// stream.
type Processor func(in Value) (out Value, keep bool, err error)

// errInvalidInput prefixes a decode failure on a single Process line.
const errInvalidInput Error = "json: invalid input"

// errDecode prefixes a streaming decode failure in Decode.
const errDecode Error = "json: decode"

// errMarshal prefixes a re-encoding failure.
const errMarshal Error = "json: marshal"

// marshal is the JSON encoder used to re-emit values. It is a package variable
// so tests can inject a failing encoder and exercise the error path. Values
// re-emitted here always originate from a successful json.Unmarshal, so the
// default json.Marshal never fails in production.
var marshal = json.Marshal

// encodeValue marshals v to compact JSON, wrapping any failure in errMarshal.
func encodeValue(v Value) ([]byte, error) {
	enc, err := marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errMarshal, err)
	}
	return enc, nil
}

// processLine decodes one JSON value from line, applies p, and re-encodes the
// result. A blank line yields no output; a dropped value (keep=false) likewise.
func processLine(p Processor, line []byte) ([][]byte, error) {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return nil, nil
	}
	var v Value
	if err := json.Unmarshal(trimmed, &v); err != nil {
		return nil, fmt.Errorf("%w %q: %w", errInvalidInput, trimmed, err)
	}
	out, keep, err := p(v)
	if err != nil || !keep {
		return nil, err
	}
	enc, err := encodeValue(out)
	if err != nil {
		return nil, err
	}
	return [][]byte{enc}, nil
}

// Process returns a command that, for each input line, decodes one JSON value,
// applies p, and re-encodes kept values as compact JSON — one value per output
// line. Blank lines are skipped. Values flow one at a time, so the command
// streams with backpressure and never buffers the whole input.
//
// This is the shared core every value-oriented json command is built on
// (see cmd-json/pluck and cmd-json/select).
func Process(p Processor) gloo.Command[[]byte, []byte] {
	return patterns.Expand(func(line []byte) ([][]byte, error) {
		return processLine(p, line)
	})
}

// decodeValues reads every JSON value framed in raw, in order. It accepts JSON
// Lines, whitespace/newline-separated values, or a pretty-printed document.
func decodeValues(raw []byte) ([]Value, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	var values []Value
	for {
		var v Value
		err := dec.Decode(&v)
		if errors.Is(err, io.EOF) {
			return values, nil
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errDecode, err)
		}
		values = append(values, v)
	}
}

// flatten streams a single top-level array element-by-element; any other value
// list passes through unchanged.
func flatten(values []Value) []Value {
	if len(values) == 1 {
		if arr, ok := values[0].([]Value); ok {
			return arr
		}
	}
	return values
}

// encodeValues marshals each value to its own compact-JSON line.
func encodeValues(values []Value) ([][]byte, error) {
	out := make([][]byte, 0, len(values))
	for _, v := range values {
		enc, err := encodeValue(v)
		if err != nil {
			return nil, err
		}
		out = append(out, enc)
	}
	return out, nil
}

// decodeLines is the buffered body of Decode: join, decode, flatten, re-encode.
func decodeLines(lines [][]byte) ([][]byte, error) {
	raw := bytes.TrimSpace(bytes.Join(lines, []byte{'\n'}))
	if len(raw) == 0 {
		return nil, nil
	}
	values, err := decodeValues(raw)
	if err != nil {
		return nil, err
	}
	return encodeValues(flatten(values))
}

// Decode normalizes arbitrary JSON framing into the one-value-per-line form the
// other json commands expect. It accepts JSON Lines, whitespace/newline
// separated values, a pretty-printed document, or a single top-level array
// (whose elements are streamed individually), and emits each value as compact
// JSON on its own line.
//
// Because reconstructing values can span line boundaries, Decode buffers its
// input; downstream value-by-value commands (Process and friends) then stream
// normally.
func Decode() gloo.Command[[]byte, []byte] {
	return patterns.Accumulate(decodeLines)
}

// AsMap returns v as a JSON object together with ok=true when v is an object;
// otherwise it returns nil, false. It is a convenience for Processors that only
// act on objects.
func AsMap(v Value) (map[string]Value, bool) {
	m, ok := v.(map[string]Value)
	return m, ok
}
