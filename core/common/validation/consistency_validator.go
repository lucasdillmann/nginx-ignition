package validation

const ValueMissingMessage = "A value is required"

type ConsistencyValidator struct {
	violations []ConsistencyViolation
}

func (v *ConsistencyValidator) Add(path string, message string) {
	violation := ConsistencyViolation{path, message}
	v.violations = append(v.violations, violation)
}

func (v *ConsistencyValidator) Result() error {
	if len(v.violations) == 0 {
		return nil
	}

	return NewError(v.violations)
}
