package format

import "github.com/tidwall/pretty"

// JSON formats the JSON body
var JSON Formatter = func(body []byte, format T) []byte {
	switch format {
	case Minified:
		return pretty.Ugly(body)
	case Expanded:
		return pretty.Pretty(body)
	}
	return body
}
