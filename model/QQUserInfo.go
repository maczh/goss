package model

import "gopkg.in/mgo.v2/bson"

type QQUserInfo struct {
	ID         bson.ObjectId `bson:"_id"`
	UserId     string        `json:"userId" bson:"userId"`
	OpenId     string        `json:"openId" bson:"openId"`         //QQ唯一账号
	NickName   string        `json:"nickName" bson:"nickName"`     //QQ昵称
	Sex        string        `json:"sex" bson:"sex"`               //QQ性别
	HeadImgUrl string        `json:"headImgUrl" bson:"headImgUrl"` //QQ头像
	BindTime   string        `json:"bindTime" bson:"bindTime"`
}
