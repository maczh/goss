package model

import "gopkg.in/mgo.v2/bson"

type AlipayUserInfo struct {
	ID         bson.ObjectId `bson:"_id"`
	UserId     string        `json:"userId" bson:"userId"`
	AlipayId   string        `json:"alipayId" bson:"alipayId"` //支付宝用户user_id
	NickName   string        `json:"nickName" bson:"nickName"` //支付宝昵称
	Sex        string        `json:"sex" bson:"sex"`           //支付宝性别
	Province   string        `json:"province" bson:"province"`
	City       string        `json:"city" bson:"city"`
	Country    string        `json:"country" bson:"country"`
	HeadImgUrl string        `json:"headImgUrl" bson:"headImgUrl"`
	BindTime   string        `json:"bindTime" bson:"bindTime"`
}
