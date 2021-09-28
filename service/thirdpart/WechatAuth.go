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

type WechatAuth struct {
	WechatUser model.WechatUserInfo
}

func NewWechatAuth() *WechatAuth {
	return &WechatAuth{}
}

func (w *WechatAuth) GetThirdUserInfo(userId, appId string) mgresult.Result {
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
	w.WechatUser, _ = mongo.GetWechatUserInfo(userId, "")
	if w.WechatUser.UserId == "" {
		return *mgresult.Error(-1, "此用户未绑定微信")
	}
	return *mgresult.Success(w.WechatUser)
}

func (w *WechatAuth) Bind(thirdUserId, userId, mobile, nickName, sex, province, city, country, headImageUrl string) mgresult.Result {
	var err error
	w.WechatUser, _ = mongo.GetWechatUserInfo(userId, thirdUserId)
	if w.WechatUser.UnionId != "" {
		//已绑定的，无需重新绑定，更新数据
		w.WechatUser.NickName = nickName
		w.WechatUser.Sex = sex
		w.WechatUser.Province = province
		w.WechatUser.City = city
		w.WechatUser.Country = country
		w.WechatUser.HeadImgUrl = headImageUrl
		mongo.UpdateWechatUserInfo(w.WechatUser)
		return *mgresult.Success(w.WechatUser)
	}
	if userId != "" {
		userInfo, _ := mysql.GetUserInfoByUserId(userId)
		if userInfo.UserId == "" {
			return *mgresult.Error(-1, "用户编号不存在")
		}
		//绑定
		w.bind(thirdUserId, userId, nickName, sex, province, city, country, headImageUrl)
		return *mgresult.Success(w.WechatUser)
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
		return *mgresult.Success(w.WechatUser)
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
	return *mgresult.Success(w.WechatUser)
}

func (w *WechatAuth) UnBind(appId, userId, thirdUserId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "无此应用代码")
	}
	w.WechatUser, _ = mongo.GetWechatUserInfo(userId, thirdUserId)
	if w.WechatUser.UnionId == "" {
		return *mgresult.Error(-1, "该微信账号未绑定")
	}
	err = mongo.DeleteWechatUserInfo(w.WechatUser)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(nil)
}

func (w *WechatAuth) Login(thirdUserId, appId, userIp, userAgent, deviceId, deviceInfo string, termType int) mgresult.Result {
	var err error
	w.WechatUser, err = mongo.GetWechatUserInfo("", thirdUserId)
	if err != nil {
		logs.Error("获取绑定数据错误:{}", err.Error())
		return *mgresult.Error(-1, err.Error())
	}
	if w.WechatUser.UnionId == "" {
		return *mgresult.Error(0, "未绑定该微信用户，请先绑定")
	}
	return service.NewUserService().LoginByUserId(w.WechatUser.UserId, appId, userIp, userAgent, deviceId, deviceInfo, termType, constant.LT_WECHAT)
}

func (w *WechatAuth) bind(thirdUserId, userId, nickName, sex, province, city, country, headImageUrl string) {
	w.WechatUser = model.WechatUserInfo{
		UserId:     userId,
		UnionId:    thirdUserId,
		NickName:   nickName,
		Sex:        sex,
		Province:   province,
		City:       city,
		Country:    country,
		HeadImgUrl: headImageUrl,
	}
	w.WechatUser, _ = mongo.InsertWechatUserInfo(w.WechatUser)
}
