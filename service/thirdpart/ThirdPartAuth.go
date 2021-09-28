package thirdpart

import "github.com/maczh/gintool/mgresult"

type ThirdPartAuth interface {
	Bind(thirdUserId, userId, mobile, nickName, sex, province, city, country, headImageUrl string) mgresult.Result
	UnBind(appId, userId, thirdUserId string) mgresult.Result
	Login(thirdUserId, appId, userIp, userAgent, deviceId, deviceInfo string, termType int) mgresult.Result
	GetThirdUserInfo(userId, appId string) mgresult.Result
}

func getThirdPartAuth(platform string) ThirdPartAuth {
	var thirdPartAuth ThirdPartAuth
	switch platform {
	case "wechat":
		thirdPartAuth = NewWechatAuth()
	case "qq":
		thirdPartAuth = NewQQAuth()
	case "alipay":
		thirdPartAuth = NewAlipayAuth()
	}
	return thirdPartAuth
}

func BindThirdUser(platform, thirdUserId, userId, mobile, nickName, sex, province, city, country, headImageUrl string) mgresult.Result {
	thirdPartAuth := getThirdPartAuth(platform)
	return thirdPartAuth.Bind(thirdUserId, userId, mobile, nickName, sex, province, city, country, headImageUrl)
}

func UnBindThirdUser(platform, appId, userId, thirdUserId string) mgresult.Result {
	thirdPartAuth := getThirdPartAuth(platform)
	return thirdPartAuth.UnBind(appId, userId, thirdUserId)
}

func LoginThirdUser(platform, thirdUserId, appId, userIp, userAgent, deviceId, deviceInfo string, termType int) mgresult.Result {
	thirdPartAuth := getThirdPartAuth(platform)
	return thirdPartAuth.Login(thirdUserId, appId, userIp, userAgent, deviceId, deviceInfo, termType)
}

func GetThirdUserInfo(platform, userId, appId string) mgresult.Result {
	thirdPartAuth := getThirdPartAuth(platform)
	return thirdPartAuth.GetThirdUserInfo(userId, appId)
}
