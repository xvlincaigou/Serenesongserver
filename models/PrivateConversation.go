package models

type Message struct {
	// 这里的 FromMe 是指消息是否是自己发送的
	FromMe  bool   `bson:"from_me"`
	Content string `bson:"content"`
}

type PrivateConversation struct {
	// 这里的 UserId 是对方的用户 ID，不是自己的用户 ID
	UserId   string    `bson:"user_id"`
	Messages []Message `bson:"messages"`
}
