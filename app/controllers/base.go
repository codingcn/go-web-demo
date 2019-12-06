package controllers

import (
	"github.com/gin-gonic/gin"
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

func (bc *BaseController) ResponseJson(c *gin.Context, httpCode int, code int, message string, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	c.JSON(httpCode, gin.H{
		"error_code":    code,
		"message": message,
		"data":    data,
	})
}
