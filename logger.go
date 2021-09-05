package log

type Logger interface {
	Log(prefix Prefix, v ...interface{})
}

type Prefix uint

const (
	Debug Prefix = iota
	Info
	Warn
	Error
)

var prefixes = []string{"DEBUG ", "INFO ", "WARN ", "ERROR "}

func (p Prefix) String() string {
	if p <= Error {
		return prefixes[p]
	}
	return "UNKNOWN "
}
