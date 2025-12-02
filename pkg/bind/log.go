package bind

type Log interface {
	Write(operation string, data any)
}

type log struct {
}

var logImpl Log = &log{}

func NewLog() Log {
	return &log{}
}
