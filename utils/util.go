package utils

import (
    "reflect"
	"log"
	"github.com/gin-gonic/gin"
)

// IsAnyFieldEmpty 检查结构体中的任何字段是否为空
func IsAnyFieldEmpty(v interface{}) bool {
    val := reflect.ValueOf(v)
    for i := 0; i < val.NumField(); i++ {
        if val.Field(i).Interface() == "" {
            return true
        }
    }
    return false
}

// handleError 统一处理错误
func HandleError(c *gin.Context, statusCode int, message string, err error) {
    c.JSON(statusCode, gin.H{"error": message})
    log.Println(err)
}