package fromtoml

import "encoding/json"

// marshal is the JSON encoder used to render the document. It is a package
// variable so tests can inject a failing encoder and exercise the error path
// (TOML decodes only into JSON-marshalable Go values, so the default never
// fails in production).
var marshal = json.Marshal
