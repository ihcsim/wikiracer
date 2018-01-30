package validator

import (
	"github.com/ihcsim/wikiracer/errors"
	"github.com/ihcsim/wikiracer/internal/wiki"
)

// InputValidator is used to validate the origin and destination inputs.
type InputValidator struct {
	wiki.Wiki
}

// NewInputValidator returns a new instance of InputValidator.
func NewInputValidator(w wiki.Wiki) *InputValidator {
	return &InputValidator{w}
}

// Validate contains rules used to validate the origin and destination inputs.
// An error is returned if either the inputs failed the rules.
func (v *InputValidator) Validate(origin, destination string) error {
	if len(origin) == 0 || len(destination) == 0 {
		return errors.InvalidEmptyInput{Origin: origin, Destination: destination}
	}

	if _, err := v.FindPage(origin); err != nil {
		return err
	}

	if _, err := v.FindPage(destination); err != nil {
		return err
	}

	return nil

}
