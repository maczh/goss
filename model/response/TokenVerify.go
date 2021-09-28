package response

type TokenVerify struct {
	Token       string `json:"token" bson:"token"`
	UserId      string `json:"userId" bson:"userId"`
	Status      int    `json:"status" bson:"status"`           //Token验证状态，1-验证通过 2-验证失败 3-token失效
	InvalidType int    `json:"invalidType" bson:"invalidType"` //失效类型 0-未失效 其他值，失效类型
	Message     string `json:"message" bson:"message"`         //验证失败与失效的错误信息
}
