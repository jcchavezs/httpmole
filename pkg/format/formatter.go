package format

// T represents the format type.
type T int

const (
	_ T = iota
	// Minified minifies the JSON response, useful for persistent logs.
	Minified
	// Expanded expands the JSON response, useful for console logs.
	Expanded
)

// Formatter is a callable type that returns a formatted body for a type.
type Formatter func(body []byte, f T) []byte

// Noop is a no-op formatter.
func Noop(body []byte, _ T) []byte {
	return body
}

// GetFormatterContentType returns the formatter based on the body content type
func GetFormatterContentType(contentType string) Formatter {
	switch contentType {
	case "application/json":
		return JSON
	default:
		return Noop
	}
}
