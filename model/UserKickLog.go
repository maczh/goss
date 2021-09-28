package model

import "gopkg.in/mgo.v2/bson"

type UserKickLog struct {
	ID          bson.ObjectId `bson:"_id"`
	Time        string        `json:"time" bson:"time"`
	RequestId   string        `json:"requestId" bson:"requestId"`
	Token       string        `json:"token" bson:"token"`             //被踢的token串
	KickedToken TokenInfo     `json:"kickedToken" bson:"kickedToken"` //被踢掉的token详情
	ByNewToken  TokenInfo     `json:"byNewToken" bson:"byNewToken"`   //新上线token
}
