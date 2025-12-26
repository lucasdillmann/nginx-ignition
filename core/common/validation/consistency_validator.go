package validation

const (
	ValueMissingMessage          = "A value is required"
	ValueCannotBeNegativeMessage = "Value must be 0 or greater"
	ValueCannotBeZeroMessage     = "Value must be greater than 0"
)

type ConsistencyValidator struct {
	violations []ConsistencyViolation
}

func NewValidator() *ConsistencyValidator {
	return &ConsistencyValidator{}
}

func (v *ConsistencyValidator) Add(path, message string) {
	violation := ConsistencyViolation{path, message}
	v.violations = append(v.violations, violation)
}

func (v *ConsistencyValidator) Result() error {
	if len(v.violations) == 0 {
		return nil
	}

	return NewError(v.violations)
}
