package garlic

import (
	"errors"
	"fmt"
)

type ExitCoder interface {
	ExitCode() int
}

type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	return fmt.Sprintf("exit %d: %s", e.Code, e.Err)
}

func (e ExitError) ExitCode() int {
	return e.Code
}

func (e ExitError) Unwrap() error {
	type unwrapper interface {
		Unwrap() error
	}
	if _, ok := e.Err.(unwrapper); ok {
		return errors.Unwrap(e.Err)
	}

	return e.Err
}

func (e ExitError) As(target any) bool {
	ex, ok := target.(*ExitError)
	if !ok {
		// if the target is not an Error,
		// pass the underlying error up the chain
		// by calling errors.As with the underlying error
		// and the target error
		return errors.As(e.Err, target)
	}

	// set the target to the current error
	(*ex) = e
	return true
}

func (e ExitError) Is(target error) bool {
	if _, ok := target.(ExitError); ok {
		// return true if target is Error
		return true
	}

	// if not, pass the underlying error up the chain
	// by calling errors.Is with the underlying error
	// and the target error
	return errors.Is(e.Err, target)
}
