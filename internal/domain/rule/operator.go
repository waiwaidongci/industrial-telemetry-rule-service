package rule

import (
	"errors"
	"fmt"
)

var ErrUnknownOperator = errors.New("unknown rule operator")

type OperatorError struct {
	Name  string
	Cause error
}

func (e *OperatorError) Error() string { return fmt.Sprintf("unknown operator %q", e.Name) }
func (e *OperatorError) Unwrap() error { return e.Cause }

type Operator func(float64, float64) bool

func ResolveOperator(name string) (Operator, error) {
	switch name {
	case ">":
		return func(a, b float64) bool { return a > b }, nil
	case ">=":
		return func(a, b float64) bool { return a >= b }, nil
	case "<":
		return func(a, b float64) bool { return a < b }, nil
	case "<=":
		return func(a, b float64) bool { return a <= b }, nil
	case "=", "==":
		return func(a, b float64) bool { return a == b }, nil
	case "!=":
		return func(a, b float64) bool { return a != b }, nil
	default:
		return nil, &OperatorError{Name: name, Cause: fmt.Errorf("%v", ErrUnknownOperator)}
	}
}
func SupportedOperators() []string       { return []string{"<", "<=", "=", "!=", ">=", ">"} }
func ValidateOperator(name string) error { _, err := ResolveOperator(name); return err }
