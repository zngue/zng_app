package validate

type Validator interface {
	Validate(req any) error
}

type NoopValidator struct{}

func (NoopValidator) Validate(req any) error {
	if validator, ok := req.(interface{ Validate() error }); ok {
		return validator.Validate()
	}
	return nil
}
