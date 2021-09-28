package model

import "gopkg.in/mgo.v2/bson"

type WechatUserInfo struct {
	ID         bson.ObjectId `bson:"_id"`
	UserId     string        `json:"userId" bson:"userId"`
	UnionId    string        `json:"unionId" bson:"unionId"`   //微信唯一账号
	NickName   string        `json:"nickName" bson:"nickName"` //微信昵称
	Sex        string        `json:"sex" bson:"sex"`           //微信性别
	Province   string        `json:"province" bson:"province"`
	City       string        `json:"city" bson:"city"`
	Country    string        `json:"country" bson:"country"`
	HeadImgUrl string        `json:"headImgUrl" bson:"headImgUrl"`
	BindTime   string        `json:"bindTime" bson:"bindTime"`
}
