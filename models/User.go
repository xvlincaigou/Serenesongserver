package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Username string `bson:"username"` // 用户名
	Avatar   string `bson:"avatar"`	// 头像保存的路径。当然这个可能之后不再需要，这个要看微信接口接的怎么样。

	// 这里的 CiRead, CiWritten, CiLiked 是用户的诗词阅读、创作、点赞记录。对于这些词，我们只保存ID，具体的诗词内容存放在Ci中
	CiRead    []primitive.ObjectID `bson:"ci_read"`
	CiWritten []primitive.ObjectID `bson:"ci_written"`
	CiLiked   []primitive.ObjectID `bson:"ci_liked"`

	Collections          []primitive.ObjectID `bson:"collections"`   // 这里我们只存放id列表，具体的收藏的诗词存放在Collections中
	SubscribedTo         []primitive.ObjectID `bson:"subscribed_to"` // 用户订阅的用户ID列表
	Subscribers          []primitive.ObjectID `bson:"subscribers"`   // 订阅本用户的用户ID列表
	Dynamics             []Dynamic            `bson:"dynamics"`		// 用户动态列表
	Drafts               []Ci                 `bson:"drafts"`		// 用户的草稿箱，这个我们选择在用户中嵌入，因为草稿箱是用户私有的
	PrivateConversations []primitive.ObjectID `bson:"private_conversations"` // 用户的私信会话列表，仍然只存放id
}
