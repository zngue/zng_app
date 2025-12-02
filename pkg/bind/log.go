package bind

type Log interface {
	Write(operation string, data any)
}

type log struct {
}

// Write implements Log.
func (l *log) Write(operation string, data any) {
	panic("unimplemented")
}

var logImpl Log = &log{}

func NewLog() Log {
	return &log{}
}
