package response

type LoginResponse struct {
	UserId   string `json:"userId" bson:"userId"`
	Token    string `json:"token" bson:"token"`
	UserName string `json:"userName" bson:"userName"`
	Mobile   string `json:"mobile"bson:"mobile"`
	Image    string `json:"image" bson:"image"`
}
