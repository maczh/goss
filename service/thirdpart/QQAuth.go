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

type QQAuth struct {
	QQUser model.QQUserInfo
}

func NewQQAuth() *QQAuth {
	return &QQAuth{}
}

func (w *QQAuth) GetThirdUserInfo(userId, appId string) mgresult.Result {
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
	w.QQUser, _ = mongo.GetQQUserInfo(userId, "")
	if w.QQUser.UserId == "" {
		return *mgresult.Error(-1, "此用户未绑定QQ")
	}
	return *mgresult.Success(w.QQUser)
}

func (w *QQAuth) Bind(thirdUserId, userId, mobile, nickName, sex, province, city, country, headImageUrl string) mgresult.Result {
	var err error
	w.QQUser, _ = mongo.GetQQUserInfo(userId, thirdUserId)
	if w.QQUser.OpenId != "" {
		//已绑定的，无需重新绑定，更新数据
		w.QQUser.NickName = nickName
		w.QQUser.Sex = sex
		w.QQUser.HeadImgUrl = headImageUrl
		mongo.UpdateQQUserInfo(w.QQUser)
		return *mgresult.Success(w.QQUser)
	}
	if userId != "" {
		userInfo, _ := mysql.GetUserInfoByUserId(userId)
		if userInfo.UserId == "" {
			return *mgresult.Error(-1, "用户编号不存在")
		}
		//绑定
		w.bind(thirdUserId, userId, nickName, sex, province, city, country, headImageUrl)
		return *mgresult.Success(w.QQUser)
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
		return *mgresult.Success(w.QQUser)
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
	return *mgresult.Success(w.QQUser)
}

func (w *QQAuth) UnBind(appId, userId, thirdUserId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "无此应用代码")
	}
	w.QQUser, _ = mongo.GetQQUserInfo(userId, thirdUserId)
	if w.QQUser.OpenId == "" {
		return *mgresult.Error(-1, "该QQ账号未绑定")
	}
	err = mongo.DeleteQQUserInfo(w.QQUser)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(nil)

}

func (w *QQAuth) Login(thirdUserId, appId, userIp, userAgent, deviceId, deviceInfo string, termType int) mgresult.Result {
	var err error
	w.QQUser, err = mongo.GetQQUserInfo("", thirdUserId)
	if err != nil {
		logs.Error("获取绑定数据错误:{}", err.Error())
		return *mgresult.Error(-1, err.Error())
	}
	if w.QQUser.OpenId == "" {
		return *mgresult.Error(0, "未绑定该QQ用户，请先绑定")
	}
	return service.NewUserService().LoginByUserId(w.QQUser.UserId, appId, userIp, userAgent, deviceId, deviceInfo, termType, constant.LT_QQ)
}

func (w *QQAuth) bind(thirdUserId, userId, nickName, sex, province, city, country, headImageUrl string) {
	w.QQUser = model.QQUserInfo{
		UserId:     userId,
		OpenId:     thirdUserId,
		NickName:   nickName,
		Sex:        sex,
		HeadImgUrl: headImageUrl,
	}
	w.QQUser, _ = mongo.InsertQQUserInfo(w.QQUser)
}
