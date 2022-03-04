package service

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/redis"
)

func (us *UserService) UpdateUserInfo(token, userId, userName, mobile, email, image, descript, realName, idCardNo string) mgresult.Result {
	tokenInfo, _ := redis.GetToken(token)
	if tokenInfo.Token == "" {
		return mgresult.Error(-1, "令牌无效")
	}
	if tokenInfo.UserId != userId {
		return mgresult.Error(-1, "令牌不匹配")
	}
	userInfo, err := mysql.GetUserInfoByUserId(userId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if mobile != "" && userInfo.Mobile != mobile {
		//更换手机号之前先检查该手机号是否存在
		u, _ := mysql.GetUserInfoByMobile(mobile)
		if u.UserId != "" {
			return mgresult.Error(-1, "要更换的手机号码已经注册，请更换成其他手机号")
		}
		userInfo.Mobile = mobile
	}
	if userName != "" {
		userInfo.UserName = userName
	}
	if email != "" {
		userInfo.Email = email
	}
	if image != "" {
		userInfo.Image = image
	}
	if descript != "" {
		userInfo.Descript = descript
	}
	if realName != "" {
		userInfo.RealName = realName
	}
	if idCardNo != "" {
		userInfo.IdCardNo = idCardNo
	}
	err = mysql.UpdateUserInfo(userInfo)
	if err != nil {
		return mgresult.Error(-1, "用户信息修改错误:"+err.Error())
	}
	return mgresult.Success(userInfo)
}

func (us *UserService) UpdateUserStatus(userId, appId string, status int) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return mgresult.Error(-1, "无此应用代码")
	}
	userInfo, err := mysql.GetUserInfoByUserId(userId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	userInfo.Status = status
	err = mysql.UpdateUserInfo(userInfo)
	if err != nil {
		return mgresult.Error(-1, "用户信息修改错误:"+err.Error())
	}
	return mgresult.Success(userInfo)
}

func (us *UserService) GetuserInfo(userId, mobile string) mgresult.Result {
	var userInfo model.UserInfo
	var err error
	if userId != "" {
		userInfo, err = mysql.GetUserInfoByUserId(userId)
	} else if mobile != "" {
		userInfo, err = mysql.GetUserInfoByMobile(mobile)
	} else {
		return mgresult.Error(-1, "用户编号和手机号不可同时为空")
	}
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return mgresult.Error(0, "无此用户")
	}
	userInfo.Password = userInfo.GetPassword()
	return mgresult.Success(userInfo)
}

func (us *UserService) ListUserTokensByApp(userId, appId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return mgresult.Error(-1, "无此应用代码")
	}
	userInfo, err := mysql.GetUserInfoByUserId(userId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return mgresult.Error(-1, "无此用户")
	}
	tokens := redis.ListTokensFromUserAppTokens(userId, appId)
	if tokens == nil || len(tokens) == 0 {
		return mgresult.Error(-1, "该用户无登录令牌")
	}
	userTokens := make([]model.TokenInfo, len(tokens))
	for i, token := range tokens {
		userTokens[i], _ = redis.GetToken(token)
	}
	return mgresult.Success(userTokens)
}
