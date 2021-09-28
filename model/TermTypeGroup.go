package model

type TermTypeGroup struct {
	GroupId   string `json:"groupId" bson:"groupId"`
	GroupName string `json:"groupName" bson:"groupName"`
	TermTypes []int  `json:"termTypes" bson:"termTypes"`
}
