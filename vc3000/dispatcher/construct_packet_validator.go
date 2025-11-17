package vc3000

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func (d *Dispatcher) validate(field interface{}, pattern string) (err error) {

	if validate == nil {
		return errors.New("internal error: validator not available")
	}

	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("panic: %s", e)
		}
	}()

	err = validate.Var(field, pattern)

	if err == nil {
		return nil
	}

	for _, e := range err.(validator.ValidationErrors) {
		if len(e.Param()) == 0 {
			return errors.Errorf("%s", e.ActualTag())
		}
		return errors.Errorf("%s %s", e.ActualTag(), e.Param())
	}

	return err

}

func init() {
	validate = validator.New()
	validate.RegisterValidation("accesspoint", ValidateAccesspoint)
}

func ValidateAccesspoint(fl validator.FieldLevel) (valid bool) {

	defer func() {
		if e := recover(); e != nil {
			valid = false
		}
	}()

	field := fl.Field().String()
	data := ([]byte(field))
	for i := range data {
		if data[i] != '0' && data[i] != '1' {
			return false
		}
	}
	return true
}
