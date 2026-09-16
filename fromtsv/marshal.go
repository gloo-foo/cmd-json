package fromtsv

import "encoding/json"

// marshal is the JSON encoder used to render rows. It is a package variable so
// tests can inject a failing encoder and exercise the error path (the default,
// json.Marshal of a string-keyed map, never fails in production).
var marshal = json.Marshal
