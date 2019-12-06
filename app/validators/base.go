package validators

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidatorError(err error) string {
	res := ""
	errs := err.(validator.ValidationErrors)

	for _, e := range errs {
		res = fmt.Sprintf("参数错误，字段%s %v", e.Field(), e.ActualTag())
	}
	return res
}
