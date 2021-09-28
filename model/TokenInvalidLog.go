package model

import "gopkg.in/mgo.v2/bson"

type TokenInvalidLog struct {
	ID          bson.ObjectId `bson:"_id"`
	Time        string        `json:"time" bson:"time"`
	Token       string        `json:"token" bson:"token"`
	TokenInfo   `json:"tokenInfo" bson:"tokenInfo"`
	CreateTime  string           `json:"createTime" bson:"createTime"`
	InvalidTime string           `json:"invalidTime" bson:"invalidTime"`
	InvalidInfo TokenInvalidInfo `json:"invalidInfo" bson:"invalidInfo"`
}

type TokenInvalidInfo struct {
	InvalidType int    `json:"invalidType" bson:"invalidType"`
	Message     string `json:"message" bson:"message"`
}
