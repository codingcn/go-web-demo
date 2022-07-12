package controllers

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/app/middlewares"
)

func init() {
	// test
	b := &BaseController{}
	b.init()
}

type BaseController struct {
}

func (bc *BaseController) init() {

}

func (bc *BaseController) GetAuthUser(ctx *gin.Context) *middlewares.UserInfoStruct {
	v, exists := ctx.Get("claims")
	if !exists {
		return nil
	}
	return &v.(*middlewares.CustomClaims).UserInfo
}
func (bc *BaseController) ResponseJson(ctx *gin.Context, httpCode int, code int, message string, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	ctx.JSON(httpCode, gin.H{
		"error_code": code,
		"message":    message,
		"data":       data,
	})
}
