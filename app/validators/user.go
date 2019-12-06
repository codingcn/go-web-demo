package validators

type BindPhoneValidate struct {
	Phone string `form:"phone" json:"phone" validate:"required"`
	Code  string `form:"error_code" json:"error_code" validate:"required"`
}