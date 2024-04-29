package val

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validator struct {
	name       string
	value      interface{}
	violations []*errdetails.BadRequest_FieldViolation
}

func NewValidator(name string, value interface{}) *Validator {
	return &Validator{
		name:  name,
		value: value,
	}
}

func (v *Validator) String() *Validator {
	if v.value == "" {
		v.violations = append(v.violations, &errdetails.BadRequest_FieldViolation{
			Field:       v.name,
			Description: fmt.Sprintf("%s is required", v.name),
		})
	}
	return v
}

func (v *Validator) MinLenght(min int) *Validator {
	value, ok := v.value.(string)
	if !ok {
		return v
	}
	if len(value) < min {
		v.violations = append(v.violations, &errdetails.BadRequest_FieldViolation{
			Field:       v.name,
			Description: fmt.Sprintf("%s must be at least %d characters", v.name, min),
		})
	}
	return v
}

func (v *Validator) MaxLenght(max int) *Validator {
	value, ok := v.value.(string)
	if !ok {
		return v
	}
	if len(value) > max {
		v.violations = append(v.violations, &errdetails.BadRequest_FieldViolation{
			Field:       v.name,
			Description: fmt.Sprintf("%s must be at most %d characters", v.name, max),
		})
	}
	return v
}

func (v *Validator) Number() *Validator {
	_, ok := v.value.(int32)

	if !ok {
		v.violations = append(v.violations, &errdetails.BadRequest_FieldViolation{
			Field:       v.name,
			Description: fmt.Sprintf("%s must be a number", v.name),
		})
	}

	return v
}

func (v *Validator) UUID() *Validator {
	_, ok := v.value.(string)
	if !ok {
		return v
	}

	_, err := uuid.Parse(v.value.(string))
	if err != nil {
		v.violations = append(v.violations, &errdetails.BadRequest_FieldViolation{
			Field:       v.name,
			Description: fmt.Sprintf("%s must be a valid UUID", v.name),
		})
	}

	return v
}

func (v *Validator) Min(min int32) *Validator {
	value, ok := v.value.(int32)
	if !ok {
		return v
	}

	if value < min {
		v.violations = append(v.violations, &errdetails.BadRequest_FieldViolation{
			Field:       v.name,
			Description: fmt.Sprintf("%s must be at least %d", v.name, min),
		})
	}
	return v
}

func (v *Validator) Validate() []*errdetails.BadRequest_FieldViolation {
	return v.violations
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusIvalid := status.New(codes.InvalidArgument, "invalid parameter")

	statusDetails, err := statusIvalid.WithDetails(badRequest)
	if err != nil {
		return statusIvalid.Err()
	}

	return statusDetails.Err()
}
