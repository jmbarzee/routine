package routine

import (
	"errors"
	"fmt"
)

var _ error = MultiError{}

// MultiError is meant to help with routines which
// produce multiple errors during their operation but
// are confined to produce a single error.
type MultiError struct {
	Errs []error
}

func (err MultiError) Error() string {
	s := ""
	for i, e := range err.Errs {
		s += fmt.Sprintf("error[%v]: %s", i, e.Error())
	}
	return s
}

func (err MultiError) As(target interface{}) bool {
	for _, e := range err.Errs {
		if errors.As(e, target) {
			return true
		}
	}
	return false
}

func (err MultiError) Is(target error) bool {
	for _, e := range err.Errs {
		if errors.Is(e, target) {
			return true
		}
	}
	return false
}
