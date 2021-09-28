package model

import "gopkg.in/mgo.v2/bson"

type UserSpecialCode struct {
	ID              bson.ObjectId `bson:"_id"`
	UserId          string        `json:"userId" bson:"userId"`
	FignerPrintCode string        `json:"fignerPrintCode,omitempty" bson:"fignerPrintCode,omitempty"` //指纹登录授权码
	FaceIdCode      string        `json:"faceIdCode,omitempty" bson:"faceIdCode,omitempty"`           //刷脸登录授权码
}
