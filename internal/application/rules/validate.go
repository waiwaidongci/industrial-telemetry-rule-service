package rules

import (
	"errors"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/rule"
)

var ErrInvalidDefinition = errors.New("invalid rule definition")

type DefinitionError struct {
	Stage string
	Cause error
}

func (e *DefinitionError) Error() string   { return fmt.Sprintf("%s: %v", e.Stage, e.Cause) }
func (e *DefinitionError) Unwrap() []error { return []error{ErrInvalidDefinition, e.Cause} }

func ValidateDefinition(d rule.Definition) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if d.Type != rule.Missing && d.Operator == "" {
		return fmt.Errorf("operator required for %s rule", d.Type)
	}
	if d.Type != rule.Missing {
		if err := rule.ValidateOperator(d.Operator); err != nil {
			return &DefinitionError{Stage: "validate operator", Cause: fmt.Errorf("%v", err)}
		}
	}
	if d.Type == rule.Rate && d.WindowSeconds < 2 {
		return fmt.Errorf("rate window must be at least two seconds")
	}
	return nil
}
func RuleDescription(d rule.Definition) string {
	return fmt.Sprintf("%s %s %s threshold %.4f over %ds", d.Name, d.Type, d.Operator, d.Threshold, d.WindowSeconds)
}
