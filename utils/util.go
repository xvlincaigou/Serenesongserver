package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

const (
	ErrMsgInvalidParams = "Invalid parameters"
	ErrMsgInternalError = "Internal server error"
	ErrMsgUserNotFound  = "User not found"
	ErrMsgMongoInsert   = "Failed to insert into MongoDB"
	ErrMsgMongoUpdate   = "Failed to update MongoDB"
	ErrMsgMongoDelete   = "Failed to delete from MongoDB"
	ErrMsgMongoFind     = "Failed to find in MongoDB"
	ErrMsgPermission    = "Permission Denied"
	ErrMsgMongoDecode   = "Failed to parse MongoDB data"
	ErrMsgMongoConnect  = "Failed to connect to MongoDB"
	ErrMsgInvalidObjID  = "Invalid object ID"
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
