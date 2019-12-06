package error_code

const (
	// 公共
	Success           = 0   // 成功
	Fail              = 1   // 失败
	ParamsCheckFailed = 400 // 校验不通过

	// 业务错误
	UserNotFound = 10000 // 参数错误
	SMSCodeError = 10001 // 参数错误
)
