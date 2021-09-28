package thirdpart

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/service"
	"github.com/maczh/logs"
)

type AlipayAuth struct {
	AlipayUser model.AlipayUserInfo
}

func NewAlipayAuth() *AlipayAuth {
	return &AlipayAuth{}
}

func (w *AlipayAuth) GetThirdUserInfo(userId, appId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "无此应用代码")
	}
	userInfo, _ := mysql.GetUserInfoByUserId(userId)
	if userInfo.UserId == "" {
		return *mgresult.Error(-1, "用户编号不存在")
	}
	w.AlipayUser, _ = mongo.GetAlipayUserInfo(userId, "")
	if w.AlipayUser.UserId == "" {
		return *mgresult.Error(-1, "此用户未绑定支付宝")
	}
	return *mgresult.Success(w.AlipayUser)
}

func (w *AlipayAuth) Bind(thirdUserId, userId, mobile, nickName, sex, province, city, country, headImageUrl string) mgresult.Result {
	var err error
	w.AlipayUser, _ = mongo.GetAlipayUserInfo(userId, thirdUserId)
	if w.AlipayUser.AlipayId != "" {
		//已绑定的，无需重新绑定，更新数据
		w.AlipayUser.NickName = nickName
		w.AlipayUser.Sex = sex
		w.AlipayUser.Province = province
		w.AlipayUser.City = city
		w.AlipayUser.Country = country
		w.AlipayUser.HeadImgUrl = headImageUrl
		mongo.UpdateAlipayUserInfo(w.AlipayUser)
		return *mgresult.Success(w.AlipayUser)
	}
	if userId != "" {
		userInfo, _ := mysql.GetUserInfoByUserId(userId)
		if userInfo.UserId == "" {
			return *mgresult.Error(-1, "用户编号不存在")
		}
		//绑定
		w.bind(thirdUserId, userId, nickName, sex, province, city, country, headImageUrl)
		return *mgresult.Success(w.AlipayUser)
	}
	if mobile != "" {
		userInfo, _ := mysql.GetUserInfoByMobile(mobile)
		if userInfo.UserId == "" {
			//完成注册新用户
			userInfo = model.UserInfo{
				UserName: nickName,
				Mobile:   mobile,
				Image:    headImageUrl,
			}
			userInfo, err = mysql.InsertUserInfo(userInfo)
			if err != nil {
				return *mgresult.Error(-1, err.Error())
			}
		}
		//绑定
		w.bind(thirdUserId, userInfo.UserId, nickName, sex, province, city, country, headImageUrl)
		return *mgresult.Success(w.AlipayUser)
	}
	//不传userId和手机号,直接用第三方用户号当作手机号注册一个新账号，然后绑定
	userInfo := model.UserInfo{
		UserName: nickName,
		Mobile:   thirdUserId,
		Image:    headImageUrl,
	}
	userInfo, err = mysql.InsertUserInfo(userInfo)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	//绑定
	w.bind(thirdUserId, userInfo.UserId, nickName, sex, province, city, country, headImageUrl)
	return *mgresult.Success(w.AlipayUser)
}

func (w *AlipayAuth) UnBind(appId, userId, thirdUserId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "无此应用代码")
	}
	w.AlipayUser, _ = mongo.GetAlipayUserInfo(userId, thirdUserId)
	if w.AlipayUser.AlipayId == "" {
		return *mgresult.Error(-1, "该支付宝账号未绑定")
	}
	err = mongo.DeleteAlipayUserInfo(w.AlipayUser)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(nil)
}

func (w *AlipayAuth) Login(thirdUserId, appId, userIp, userAgent, deviceId, deviceInfo string, termType int) mgresult.Result {
	var err error
	w.AlipayUser, err = mongo.GetAlipayUserInfo("", thirdUserId)
	if err != nil {
		logs.Error("获取绑定数据错误:{}", err.Error())
		return *mgresult.Error(-1, err.Error())
	}
	if w.AlipayUser.AlipayId == "" {
		return *mgresult.Error(0, "未绑定该支付宝用户，请先绑定")
	}
	return service.NewUserService().LoginByUserId(w.AlipayUser.UserId, appId, userIp, userAgent, deviceId, deviceInfo, termType, constant.LT_ALIPAY)
}

func (w *AlipayAuth) bind(thirdUserId, userId, nickName, sex, province, city, country, headImageUrl string) {
	w.AlipayUser = model.AlipayUserInfo{
		UserId:     userId,
		AlipayId:   thirdUserId,
		NickName:   nickName,
		Sex:        sex,
		Province:   province,
		City:       city,
		Country:    country,
		HeadImgUrl: headImageUrl,
	}
	w.AlipayUser, _ = mongo.InsertAlipayUserInfo(w.AlipayUser)
}
