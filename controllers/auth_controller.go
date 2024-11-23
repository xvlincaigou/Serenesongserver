package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	// 获取请求参数 wxcode
	wxcode := c.Query("wxcode")
	if wxcode == "" {
		// 如果 wxcode 为空，返回错误响应
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}

	services.LoginHandler(c, wxcode)
}
