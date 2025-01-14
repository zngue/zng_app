package validate

type Validator interface {
	Validate() error
}

func Validate(v any) (err error) {
	if rs, ok := v.(Validator); ok {
		return rs.Validate()
	}
	return
}
