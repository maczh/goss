package service

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/nacos"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"time"
)

const (
	KEY_USER_REG_SMS_CODE   = "user:sms:code:register:"
	KEY_USER_LOGIN_SMS_CODE = "user:sms:code:login:"
)

type UserService struct {
	User model.UserInfo
}

func NewUserService() *UserService {
	return &UserService{}
}

//发送注册短信
func (us *UserService) SendRegisterSms(appId, mobile string) mgresult.Result {
	var err error
	if mobile == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	us.User, err = mysql.GetUserInfoByMobile(mobile)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if us.User.UserId != "" {
		return mgresult.Error(-1, "该手机号已经注册，请直接登录")
	}
	smsCode := utils.GetRandomIntString(6)
	appSettings, err := mongo.GetAppSettings(appId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if appSettings.AppId == "" {
		return mgresult.Error(-1, "应用代码不正确")
	}
	template := appSettings.SmsRegisterTemplate
	signcode := appSettings.SmsSignCode
	json := map[string]string{"code": smsCode}
	err = nacos.SendSms(mobile, template, signcode, utils.ToJSON(json))
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return mgresult.Error(-1, "系统异常")
	}
	goredis.Set(KEY_USER_REG_SMS_CODE+mobile, smsCode, 1*time.Minute)
	return mgresult.Success(nil)
}

//用户注册短信确认
func (us *UserService) RegisterBySmsCode(mobile, smsCode, userName, email, image, descript, appId, deviceId, userIp, userAgent, deviceInfo string, termType int) mgresult.Result {
	var err error
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return mgresult.Error(-1, "系统异常")
	}
	code := goredis.Get(KEY_USER_REG_SMS_CODE + mobile).Val()
	if code == "" {
		return mgresult.Error(-1, "短信验证码已过期，请重新发送")
	}
	if smsCode != code {
		return mgresult.Error(-1, "短信验证码错误")
	}
	//完成注册新用户
	if userName == "" {
		userName = "手机用户" + utils.Right(mobile, 4)
	}
	us.User = model.UserInfo{
		UserName: userName,
		Mobile:   mobile,
		Email:    email,
		Image:    image,
		Descript: descript,
	}
	us.User, err = mysql.InsertUserInfo(us.User)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	result, err := us.login(appId, deviceId, userIp, userAgent, deviceInfo, termType, constant.LT_REGISTER_SMS)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return mgresult.Success(result)
}

//手机号+密码注册
func (us *UserService) RegisterByPassword(mobile, password, userName, email, image, descript, appId, deviceId, userIp, userAgent, deviceInfo string, termType int) mgresult.Result {
	var err error
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	if password == "" || len(password) < 6 {
		return mgresult.Error(-1, "密码不可为空或密码太短，最少6位")
	}
	//完成注册新用户
	if userName == "" {
		userName = "手机用户" + utils.Right(mobile, 4)
	}
	us.User = model.UserInfo{
		UserName: userName,
		Mobile:   mobile,
		Password: utils.Base64Decode(password),
		Email:    email,
		Image:    image,
		Descript: descript,
	}
	us.User, err = mysql.InsertUserInfo(us.User)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	result, err := us.login(appId, deviceId, userIp, userAgent, deviceInfo, termType, constant.LT_REGISTER_OTHER)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return mgresult.Success(result)
}

func (us *UserService) UserTransfer(userId, userName, password, mobile, image, createTime string, status int) mgresult.Result {
	var err error
	if !utils.IsChinaMobileString(mobile) {
		return mgresult.Error(-1, "手机号格式不正确")
	}
	if status == -1 {
		status = constant.DISABLED
	}
	//完成注册新用户
	us.User = model.UserInfo{
		UserId:     userId,
		UserName:   userName,
		Mobile:     mobile,
		Image:      image,
		Password:   password,
		CreateTime: createTime,
		Status:     status,
	}
	us.User, err = mysql.InsertUserInfo(us.User)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	return mgresult.Success(us.User)
}
