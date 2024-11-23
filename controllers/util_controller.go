package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCiById(c *gin.Context) {
	// 获取请求参数 _id
	_id := c.Query("_id")
	if _id == "" {
		// 如果 _id 为空，返回错误响应
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}

	services.GetCiByIdHandler(c, _id)
}
