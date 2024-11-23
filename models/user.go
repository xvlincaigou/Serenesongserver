package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"`         // 用户的ID，使用mongodb的ObjectId作为ID
	OpenID     string             `bson:"openid"`      // 用户的OpenID，这是微信返回的、用户在一个小程序中的唯一标识
	SessionKey string             `bson:"session_key"` // 用户的SessionKey，这是微信返回的会话密钥
	Token      string             `bson:"token"`       // 用户的Token，这是我们自己生成的、用于用户登录验证的令牌

	Avatar string `bson:"avatar"` // 头像保存的路径。当然这个可能之后不再需要，这个要看微信接口接的怎么样。

	// 这里的 CiWritten 是用户的诗词创作记录。对于这些词，我们只保存ID，具体的诗词内容存放在Ci中
	CiWritten []primitive.ObjectID `bson:"ci_written"`

	Collections          []primitive.ObjectID `bson:"collections"`           // 这里我们只存放id列表，具体的收藏的诗词存放在Collections中
	SubscribedTo         []primitive.ObjectID `bson:"subscribed_to"`         // 用户订阅的用户ID列表
	Subscribers          []primitive.ObjectID `bson:"subscribers"`           // 订阅本用户的用户ID列表
	Dynamics             []Dynamic            `bson:"dynamics"`              // 用户动态列表
	Drafts               []primitive.ObjectID `bson:"drafts"`                // 用户的草稿箱，这个我们选择在用户中嵌入，因为草稿箱是用户私有的
	PrivateConversations []primitive.ObjectID `bson:"private_conversations"` // 用户的私信会话列表，仍然只存放id
}

// NewUser 创建一个初始化后的 User 实例，确保切片字段不为 nil
func NewUser(openID, sessionKey, token string) User {
	return User{
		OpenID:               openID,
		SessionKey:           sessionKey,
		Token:                token,
		Avatar:               "",
		CiWritten:            []primitive.ObjectID{},
		Collections:          []primitive.ObjectID{},
		SubscribedTo:         []primitive.ObjectID{},
		Subscribers:          []primitive.ObjectID{},
		Dynamics:             []Dynamic{},
		Drafts:               []primitive.ObjectID{},
		PrivateConversations: []primitive.ObjectID{},
	}
}
