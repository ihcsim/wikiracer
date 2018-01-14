package wikiracer

// Validator can perform validations.
type Validator interface {

	// Validate contains rules used to validate the origin and destination inputs.
	// An error is returned if either the inputs don't comply with the rules.
	Validate(origin, destination string) error
}
