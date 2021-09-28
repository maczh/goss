package model

import (
	"github.com/maczh/gintool/mgresult"
	"gopkg.in/mgo.v2/bson"
)

type UserLogoutLog struct {
	ID        bson.ObjectId     `bson:"_id"`
	Time      string            `json:"time" bson:"time"`
	RequestId string            `json:"requestId" bson:"requestId"`
	UserId    string            `json:"userId" bson:"userId"`
	AppId     string            `json:"appId" bson:"appId"`
	UserAgent string            `json:"userAgent" bson:"userAgent"`
	TermType  int               `json:"termType" bson:"termType"`
	Request   map[string]string `json:"request" bson:"request"`
	Response  mgresult.Result   `json:"response" bson:"response"`
}
