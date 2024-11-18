package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

const (
	ErrMsgInvalidParams = "invalid parameters"
	ErrMsgInternalError = "internal server error"
	ErrMsgUserNotFound  = "user not found"
	ErrMsgMongoInsert   = "failed to insert into MongoDB"
	ErrMsgMongoUpdate   = "failed to update MongoDB"
	ErrMsgMongoDelete   = "failed to delete from MongoDB"
	ErrMsgMongoFind     = "failed to find in MongoDB"
	ErrMsgPermission    = "permission denied"
	ErrMsgMongoDecode   = "failed to parse MongoDB data"
	ErrMsgInvalidObjID  = "invalid object ID"
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
	if err != nil {
		message = message + ": " + err.Error()
	}
	c.JSON(statusCode, gin.H{"message": message})
	log.Println(err)
}
