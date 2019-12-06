package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-demo/app/constants"
	"go-web-demo/app/constants/error_code"
	"go-web-demo/app/controllers"
	"go-web-demo/app/middlewares"
	"go-web-demo/app/repositories"
	"go-web-demo/app/validators"
	"go-web-demo/kernel/goredis"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type UserController struct {
	controllers.BaseController
}

// @Tags 用户API
// @Router /user/info/bind_phone [post]
// @Security Bearer
// @Summary 绑定手机号
// @Description 绑定手机号
// @Accept  json
// @Produce  json
// @Param phone formData string true "手机号"
// @Param code formData string true "验证码"
// @Success 200 {object} controllers.HTTPSuccess
// @Failure 400 {object} controllers.HTTPError

func (u UserController) BindPhone(ctx *gin.Context) {
	// 安全校验
	var params validators.BindPhoneValidate
	err := ctx.ShouldBind(&params)
	if err != nil {
		u.ResponseJson(ctx, http.StatusBadRequest, error_code.ParamsCheckFailed, "参数错误", err)
		return
	}

	err = validator.New().Struct(params)
	if err != nil {
		u.ResponseJson(ctx, http.StatusBadRequest, error_code.ParamsCheckFailed, "参数错误", validators.NewValidatorError(err))
		return
	}
	// 从上下文中获取jwt认证用户信息
	userInfo := ctx.MustGet("claims").(*middlewares.CustomClaims).UserInfo

	redisDefault := goredis.Connect.Default
	code := redisDefault.Get(fmt.Sprintf(constants.SMSSendCodeKey, userInfo.UserId, params.Phone)).Val()
	if code != params.Code {
		u.ResponseJson(ctx, http.StatusBadRequest, error_code.SMSCodeError, "验证码错误", nil)
		return
	}

	err = repositories.UpdateUserPhoneByUserId(userInfo.UserId, params.Phone)
	if err != nil {
		u.ResponseJson(ctx, http.StatusOK, error_code.Success, "绑定失败", nil)
	}
	u.ResponseJson(ctx, http.StatusOK, error_code.Success, "绑定成功", nil)
	// 验证成功，丢弃验证码
	go redisDefault.Del(fmt.Sprintf(constants.SMSSendCodeKey, userInfo.UserId, params.Phone))
}
