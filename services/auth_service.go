package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	loginURL = "https://api.weixin.qq.com/sns/jscode2session"
)

// LoginResponse 定义了微信登录接口的响应结构
type LoginResponse struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// LoginHandler 处理微信登录请求
func LoginHandler(c *gin.Context, wxcode string) {
	// 创建 HTTP 客户端
	client := &http.Client{}
	// 构建请求 URL
	url := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", loginURL, config.AppID, config.Secret, wxcode)
	// 创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// 如果创建请求失败，返回错误响应
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	}

	// 发送 HTTP 请求
	res, err := client.Do(req)
	if err != nil {
		// 如果发送请求失败，返回错误响应
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	}
	defer res.Body.Close()

	// 解析响应体
	var loginResp LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&loginResp); err != nil {
		// 如果解析响应体失败，返回错误响应
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	}
	// 检查响应中的错误码
	if loginResp.ErrCode != 0 {
		// 如果有错误码，返回错误响应
		utils.HandleError(c, http.StatusInternalServerError, loginResp.ErrMsg, nil)
		return
	}

	// 生成 token
	hash := md5.New()
	hash.Write([]byte(loginResp.SessionKey + loginResp.OpenId))
	token := hex.EncodeToString(hash.Sum(nil))

	// 检查数据库中是否有相同 openid 的用户
	collection := config.MongoClient.Database("serenesong").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"openid": loginResp.OpenId}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		// 如果没有找到用户，则创建新用户
		var n_user = models.NewUser(loginResp.OpenId, loginResp.SessionKey, token)
		_, err = collection.InsertOne(ctx, n_user)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
			return
		}
	} else if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	} else {
		// 如果找到用户，则更新 token 和 session_key
		update := bson.M{
			"$set": bson.M{
				"session_key": loginResp.SessionKey,
				"token":       token,
			},
		}
		_, err = collection.UpdateOne(ctx, bson.M{"openid": loginResp.OpenId}, update)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
			return
		}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"token": token})

}
