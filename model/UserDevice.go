package model

import "gopkg.in/mgo.v2/bson"

type UserDevice struct {
	ID         bson.ObjectId `bson:"_id"`
	DeviceId   string        `json:"deviceId" bson:"deviceId"`
	UserId     string        `json:"userId" bson:"userId"`
	DeviceInfo string        `json:"deviceInfo" bson:"deviceInfo"`
}
