package service

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/logic"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/model/response"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/nacos"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgtrace"
	"github.com/maczh/utils"
	"time"
)

func (us *UserService) loadUser(userId, mobile string) *UserService {
	if userId != "" {
		us.User, _ = mysql.GetUserInfoByUserId(userId)
	} else if mobile != "" {
		us.User, _ = mysql.GetUserInfoByMobile(mobile)
	}
	return us
}

func (us *UserService) login(appId, deviceId, userIp, userAgent, deviceInfo string, termType, loginType int) (response.LoginResponse, error) {
	//自动登录，生成token
	tokenInfo, err := logic.NewUserToken(us.User.UserId, appId, deviceId, userAgent, userIp, deviceInfo, termType)
	if err != nil {
		return response.LoginResponse{}, err
	}
	//生成返回的结果数据
	result := response.LoginResponse{
		Token:    tokenInfo.Token,
		UserId:   us.User.UserId,
		Mobile:   us.User.Mobile,
		UserName: us.User.UserName,
		Image:    us.User.Image,
	}
	//记录登录日志
	loginLog := model.UserLoginLog{
		UserId:    us.User.UserId,
		RequestId: mgtrace.GetRequestId(),
		AppId:     appId,
		UserAgent: userAgent,
		UserIp:    userIp,
		TermType:  termType,
		LoginType: loginType,
	}
	_, err = mongo.InsertLoginLog(loginLog)
	if err != nil {
		logs.Error("保存登录日志错误：{}", err.Error())
	}
	return result, nil
}

func (us *UserService) SendLoginSms(appId, mobile string) mgresult.Result {
	var err error
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return mgresult.Error(-1, "系统异常")
	}
	us.User, err = mysql.GetUserInfoByMobile(mobile)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if us.User.UserId == "" {
		return mgresult.Error(-1, "此手机号未注册，请先注册")
	}
	smsCode := utils.GetRandomIntString(6)
	appSettings, err := mongo.GetAppSettings(appId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if appSettings.AppId == "" {
		return mgresult.Error(-1, "应用代码不正确")
	}
	template := appSettings.SmsLoginTemplate
	signcode := appSettings.SmsSignCode
	json := map[string]string{"code": smsCode}
	err = nacos.SendSms(mobile, template, signcode, utils.ToJSON(json))
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	goredis.Set(KEY_USER_LOGIN_SMS_CODE+mobile, smsCode, 1*time.Minute)
	return mgresult.Success(nil)
}

//用户短信确认登录
func (us *UserService) LoginBySmsCode(mobile, smsCode, appId, deviceId, userIp, userAgent, deviceInfo string, termType int) mgresult.Result {
	var err error
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return mgresult.Error(-1, "系统异常")
	}
	code := goredis.Get(KEY_USER_LOGIN_SMS_CODE + mobile).Val()
	if code == "" {
		return mgresult.Error(-1, "短信验证码已过期，请重新发送")
	}
	if smsCode != code {
		return mgresult.Error(-1, "短信验证码错误")
	}
	userInfo, err := mysql.GetUserInfoByMobile(mobile)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return mgresult.Error(-1, "该手机号用户未注册")
	}
	if userInfo.Status != constant.NORMAL {
		return mgresult.Error(-1, "用户状态不可用")
	}
	b, err := NewAppService().CheckAppTermType(appId, termType)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if !b {
		return mgresult.Error(-1, "本应用不支持此终端类型登录")
	}
	us.loadUser("", mobile)
	result, err := us.login(appId, deviceId, userIp, userAgent, deviceInfo, termType, constant.LT_SMS)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return mgresult.Success(result)
}

//用户密码登录
func (us *UserService) LoginByPassword(mobile, password, appId, deviceId, userIp, userAgent, deviceInfo string, termType int) mgresult.Result {
	var err error
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	userInfo, err := mysql.GetUserInfoByMobile(mobile)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return mgresult.Error(-1, "该手机号用户未注册")
	}
	if userInfo.Status != constant.NORMAL {
		return mgresult.Error(-1, "用户状态不可用")
	}
	if utils.Base64Decode(password) != userInfo.GetPassword() {
		return mgresult.Error(-1, "用户密码错误")
	}
	b, err := NewAppService().CheckAppTermType(appId, termType)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if !b {
		return mgresult.Error(-1, "本应用不支持此终端类型登录")
	}
	us.loadUser("", mobile)
	result, err := us.login(appId, deviceId, userIp, userAgent, deviceInfo, termType, constant.LT_USERPWD)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return mgresult.Success(result)
}

//手机号一键登录
func (us *UserService) LoginByAliMobile(mobileToken, appId, deviceId, userIp, userAgent, deviceInfo string, termType int) mgresult.Result {
	mobile, err := nacos.GetMobileNumber(utils.GetUUIDString(), mobileToken)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	userInfo, err := mysql.GetUserInfoByMobile(mobile)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return mgresult.Success(response.LoginResponse{Mobile: mobile})
	}
	if userInfo.Status != constant.NORMAL {
		return mgresult.Error(-1, "用户状态不可用")
	}
	b, err := NewAppService().CheckAppTermType(appId, termType)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if !b {
		return mgresult.Error(-1, "本应用不支持此终端类型登录")
	}
	result, err := us.loadUser("", mobile).login(appId, deviceId, userIp, userAgent, deviceInfo, termType, constant.LT_MOBILE_ONEKEY)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return mgresult.Success(result)
}

func (us *UserService) LoginByUserId(userId, appId, deviceId, userIp, userAgent, deviceInfo string, termType, loginType int) mgresult.Result {
	userInfo, err := mysql.GetUserInfoByUserId(userId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return mgresult.Error(-1, "该手机号用户未注册")
	}
	if userInfo.Status != constant.NORMAL {
		return mgresult.Error(-1, "用户状态不可用")
	}
	b, err := NewAppService().CheckAppTermType(appId, termType)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if !b {
		return mgresult.Error(-1, "本应用不支持此终端类型登录")
	}
	us.loadUser(userId, "")
	result, err := us.login(appId, deviceId, userIp, userAgent, deviceInfo, termType, loginType)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return mgresult.Success(result)
}
