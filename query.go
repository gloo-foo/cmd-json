package command

import (
	"bytes"
	"context"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
	"github.com/gomatic/cirql"
)

// QueryScript is the source text of a cirql pipeline query.
type QueryScript string

// Query returns a command that runs a cirql pipeline query over the JSON input.
// The full input is decoded into a result set (accepting JSON Lines, a stream of
// values, or a single top-level array), run through the cirql pipeline, and each
// result emitted as a compact JSON line. An invalid query fails the stream.
//
// cirql is the GraphQL-influenced JSON pipeline language (map/filter/reduce/
// sort/flatMap/limit/uniq over .field expressions) — this is how cmd-json rivals
// jq in pure Go. See github.com/gomatic/cirql.
//
//	cmd-json query 'filter .stars > 1000 | map { name: .name } | sort .name | limit 10'
func Query(q QueryScript) gloo.Command[[]byte, []byte] {
	pipeline, err := cirql.Parse(cirql.Query(q))
	if err != nil {
		return errorCommand(ErrQuery.With(err))
	}
	return patterns.Accumulate(func(lines [][]byte) ([][]byte, error) {
		return runQuery(pipeline, lines)
	})
}

// runQuery decodes the accumulated input, runs the pipeline, and encodes each
// result value as a compact-JSON line.
func runQuery(pipeline cirql.Pipeline, lines [][]byte) ([][]byte, error) {
	raw := bytes.TrimSpace(bytes.Join(lines, []byte{'\n'}))
	if len(raw) == 0 {
		return nil, nil
	}
	values, err := decodeValues(raw)
	if err != nil {
		return nil, err
	}
	out, err := pipeline.Run(flatten(values))
	if err != nil {
		return nil, ErrQuery.With(err)
	}
	return encodeValues(out)
}

// errorCommand returns a command that ignores its input and fails with err.
func errorCommand(err error) gloo.Command[[]byte, []byte] {
	return gloo.FuncCommand[[]byte, []byte](func(ctx context.Context, _ gloo.Stream[[]byte]) gloo.Stream[[]byte] {
		return gloo.Generate(ctx, func(_ context.Context, _ func([]byte) bool, sendErr func(error)) {
			sendErr(err)
		})
	})
}
