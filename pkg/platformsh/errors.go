package platformsh

import (
	"fmt"
	"strings"
)

type AggregateError []error

//noinspection GoAssignmentToReceiver
func (e AggregateError) Append(errors ...error) AggregateError {
	for _, err := range errors {
		switch err := err.(type) {
		case AggregateError:
			e = append(e, err...)
		default:
			e = append(e, err)
		}
	}
	return e
}

func (e AggregateError) Error() string {
	errors := make([]string, len(e))
	for idx, err := range e {
		errors[idx] = err.Error()
	}

	return strings.Join(errors, ", ")
}

type MissingEnvironment struct {
	Name       string
	InnerError error
}

func (e MissingEnvironment) Error() string {
	s := fmt.Sprintf("no environment variable found for %s", e.Name)
	if e.InnerError != nil {
		s = fmt.Sprintf("%s: %v", s, e.InnerError)
	}
	return s
}

func missingEnvironment(names ...string) error {
	if len(names) == 0 {
		return nil
	}

	if len(names) == 1 {
		return MissingEnvironment{Name: names[0]}
	}

	agg := make(AggregateError, len(names))
	for idx, name := range names {
		agg[idx] = MissingEnvironment{Name: name}
	}

	return agg
}
