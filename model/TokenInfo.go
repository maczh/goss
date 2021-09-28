package model

type TokenInfo struct {
	Token      string `json:"token" bson:"token"`
	UserId     string `json:"userId" bson:"userId"`
	AppId      string `json:"appId" bson:"appId"`
	DeviceId   string `json:"deviceId" bson:"deviceId"`     //设备ID
	UserAgent  string `json:"userAgent" bson:"userAgent"`   //当前的
	UserIp     string `json:"userIp" bson:"userIp"`         //当前的用户终端IP地址
	TermType   int    `json:"termType" bson:"termType"`     //终端分类ID
	Status     int    `json:"status" bson:"status"`         //状态
	Error      string `json:"error" bson:"error"`           //失效信息
	CreateTime string `json:"createTime" bson:"createTime"` //生成时间
}
