package model

import "gopkg.in/mgo.v2/bson"

type AppSettings struct {
	ID                  bson.ObjectId `bson:"_id"`
	AppId               string        `json:"appId" bson:"appId"`
	TokenTtl            int           `json:"tokenTtl" bson:"tokenTtl"`                       //token失效时间设置，单位为分钟
	TermTypes           []int         `json:"termTypes" bson:"termTypes"`                     //支持的终端类型
	MaxOnlineTokens     int           `json:"maxOnlineTokens" bson:"maxOnlineTokens"`         //每用户保留最大token数,0为不限
	SmsSignCode         string        `json:"smsSignCode" bson:"smsSignCode"`                 //应用短信签名
	SmsLoginTemplate    string        `json:"smsLoginTemplate" bson:"smsLoginTemplate"`       //应用登录短信模板代码
	SmsRegisterTemplate string        `json:"smsRegisterTemplate" bson:"smsRegisterTemplate"` //应用注册用户短信模板代码
	VerifySign          bool          `json:"verifySign" bson:"verifySign"`                   //该应用是否需要验证签名
}
