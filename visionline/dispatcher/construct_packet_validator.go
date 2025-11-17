package visionline

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func (d *Dispatcher) validate(field interface{}, pattern string) (err error) {

	if validate == nil {
		return errors.New("internal error: validator not available")
	}

	if len(pattern) == 0 {
		return nil
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
	validate.RegisterValidation("phone", ValidatePhone)
}

func ValidatePhone(fl validator.FieldLevel) (valid bool) {

	defer func() {
		if e := recover(); e != nil {
			valid = false
		}
	}()

	field := fl.Field().String()
	field = strings.Replace(field, " ", "", -1)
	field = strings.TrimLeft(field, "+0")

	if len(field) == 0 {
		return false
	}

	if number := cast.ToUint64(field); number == 0 {
		return false
	}

	return true
}
