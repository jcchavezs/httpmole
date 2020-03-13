package format

type T int

const (
	_ T = iota
	// Minified minifies the JSON response, useful for persistent logs.
	Minified
	// Expanded expands the JSON response, useful for console logs.
	Expanded
)

type Formatter func(body []byte, f T) []byte

func Noop(body []byte, _ T) []byte {
	return body
}

func GetFormatterContentType(contentType string) Formatter {
	switch contentType {
	case "application/json":
		return JSON
	default:
		return Noop
	}
}
