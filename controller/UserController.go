package controller

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/service"
	"github.com/maczh/utils"
	"strconv"
)

// SendRegisterSms	godoc
// @Summary		注册新用户，发送短信验证码
// @Description	注册新用户，发送短信验证码
// @Tags	用户注册
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/register/sms [post][get]
func SendRegisterSms(params map[string]string) mgresult.Result {
	return service.NewUserService().SendRegisterSms(params["appId"], params["mobile"])
}

// RegisterBySmsCode	godoc
// @Summary		短信验证码注册
// @Description	短信验证码注册新用户
// @Tags	用户注册
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	smsCode formData string true "短信验证码"
// @Param	userName formData string false "用户昵称"
// @Param	email formData string false "email地址"
// @Param	image formData string false "头像图片url"
// @Param	descript formData string false "用户个性化简介"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用，可以由应用方自定义"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/register/sms/confirm [post][get]
func RegisterBySmsCode(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobile") || params["mobile"] == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.Exists(params, "smsCode") || params["smsCode"] == "" {
		return mgresult.Error(-1, "短信验证码不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewUserService().RegisterBySmsCode(params["mobile"], params["smsCode"], params["userName"], params["email"], params["image"], params["descript"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// RegisterByPassword	godoc
// @Summary		用密码方式注册新用户
// @Description	用密码方式注册新用户
// @Tags	用户注册
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	password formData string true "用户密码，用Base64编码"
// @Param	userName formData string false "用户昵称"
// @Param	email formData string false "email地址"
// @Param	image formData string false "头像图片url"
// @Param	descript formData string false "用户个性化简介"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/register/pwd [post][get]
func RegisterByPassword(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobile") || params["mobile"] == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.Exists(params, "password") || params["password"] == "" {
		return mgresult.Error(-1, "密码不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewUserService().RegisterByPassword(params["mobile"], params["password"], params["userName"], params["email"], params["image"], params["descript"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// SendLoginSms	godoc
// @Summary		发送登录短信验证码
// @Description	发送登录短信验证码
// @Tags	用户登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/sms [post][get]
func SendLoginSms(params map[string]string) mgresult.Result {
	return service.NewUserService().SendLoginSms(params["appId"], params["mobile"])
}

// LoginBySmsCode	godoc
// @Summary		短信验证码登录
// @Description	短信验证码登录
// @Tags	用户登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	smsCode formData string true "短信验证码"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/sms/confirm [post][get]
func LoginBySmsCode(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobile") || params["mobile"] == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.Exists(params, "smsCode") || params["smsCode"] == "" {
		return mgresult.Error(-1, "短信验证码不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewUserService().LoginBySmsCode(params["mobile"], params["smsCode"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// LoginByPassword	godoc
// @Summary		用手机号+密码登录
// @Description	用手机号+密码方式登录
// @Tags	用户登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	password formData string true "用户密码,用Base64编码"
// @Param	appId formData string true "应用编号"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/pwd [post][get]
func LoginByPassword(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobile") || params["mobile"] == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.Exists(params, "password") || params["password"] == "" {
		return mgresult.Error(-1, "密码不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewUserService().LoginByPassword(params["mobile"], params["password"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// GetuserInfo	godoc
// @Summary		获取用户详情
// @Description	获取用户详情
// @Tags	用户管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string false "用户编号"
// @Param	mobile formData string false "用户手机号码"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/user/get [post][get]
func GetuserInfo(params map[string]string) mgresult.Result {
	return service.NewUserService().GetuserInfo(params["userId"], params["mobile"])
}

// UpdateUserStatus	godoc
// @Summary		修改账户状态
// @Description	修改账户状态
// @Tags	用户管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户编号"
// @Param	status formData int true "账户状态 1-正常 2-禁用 3-无效"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/user/status/update [post][get]
func UpdateUserStatus(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "userId") || params["userId"] == "" {
		return mgresult.Error(-1, "用户编号不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "status") || params["status"] == "" {
		return mgresult.Error(-1, "状态值不可为空")
	}
	status, _ := strconv.Atoi(params["status"])
	return service.NewUserService().UpdateUserStatus(params["userId"], params["appId"], status)
}

// UpdateUserInfo	godoc
// @Summary		修改用户详情
// @Description	修改用户详情
// @Tags	用户管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户编号"
// @Param	token formData string true "用户令牌token"
// @Param	userName formData string false "用户昵称"
// @Param	mobile formData string false "手机号，唯一"
// @Param	email formData string false "email地址"
// @Param	image formData string false "头像图片url"
// @Param	descript formData string false "用户个性化简介"
// @Param	realName formData string false "用户实名"
// @Param	idCardNo formData string false "用户身份证号"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/user/update [post][get]
func UpdateUserInfo(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "userId") || params["userId"] == "" {
		return mgresult.Error(-1, "用户编号不可为空")
	}
	if !utils.Exists(params, "token") || params["token"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewUserService().UpdateUserInfo(params["token"], params["userId"], params["userName"], params["mobile"], params["email"], params["image"], params["descript"], params["realName"], params["idCardNo"])
}

// GetFingerPrintCode	godoc
// @Summary		获取用户指纹登录码，必须在用户登录状态
// @Description	获取用户指纹登录码，必须在用户登录状态
// @Tags	指纹登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户编号"
// @Param	token formData string true "用户令牌token"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/finger/auth [post][get]
func GetFingerPrintCode(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "userId") || params["userId"] == "" {
		return mgresult.Error(-1, "用户编号不可为空")
	}
	if !utils.Exists(params, "token") || params["token"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewFingerPrintService().GetFingerPrintCode(params["userId"], params["token"])
}

// LoginByFingerPrintCode	godoc
// @Summary		用手机号+指纹码登录
// @Description	用手机号+指纹码方式登录
// @Tags	指纹登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	fingerPrintCode formData string true "用户指纹码"
// @Param	appId formData string true "应用编号"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/finger [post][get]
func LoginByFingerPrintCode(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobile") || params["mobile"] == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.Exists(params, "fingerPrintCode") || params["fingerPrintCode"] == "" {
		return mgresult.Error(-1, "指纹码不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewFingerPrintService().LoginByFingerPrintCode(params["mobile"], params["fingerPrintCode"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// GetFaceIdCode	godoc
// @Summary		获取用户faceid登录码，必须在用户登录状态
// @Description	获取用户faceid登录码，必须在用户登录状态
// @Tags	扫脸登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户编号"
// @Param	token formData string true "用户令牌token"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/faceid/auth [post][get]
func GetFaceIdCode(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "userId") || params["userId"] == "" {
		return mgresult.Error(-1, "用户编号不可为空")
	}
	if !utils.Exists(params, "token") || params["token"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewFaceIdService().GetFaceIdCode(params["userId"], params["token"])
}

// LoginByFaceIdCode	godoc
// @Summary		用手机号+扫脸登录
// @Description	用手机号+扫脸方式登录
// @Tags	扫脸登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobile formData string true "用户手机号码"
// @Param	faceIdCode formData string true "用户脸部识别码"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/faceid [post][get]
func LoginByFaceIdCode(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobile") || params["mobile"] == "" {
		return mgresult.Error(-1, "手机号不可为空")
	}
	if !utils.Exists(params, "faceIdCode") || params["faceIdCode"] == "" {
		return mgresult.Error(-1, "脸部识别码不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewFaceIdService().LoginByFaceIdCode(params["mobile"], params["faceIdCode"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// LoginByAliMobile	godoc
// @Summary		用手机号一键登录
// @Description	用手机号一键登录
// @Tags	手机号一键登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	mobileToken formData string true "用户手机号码验证token,来自阿里云sdk生成的"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/mobile [post][get]
func LoginByAliMobile(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "mobileToken") || params["mobileToken"] == "" {
		return mgresult.Error(-1, "一键登录token不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewUserService().LoginByAliMobile(params["mobileToken"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// UserTransfer	godoc
// @Summary		用户数据逐条导入
// @Description	用户数据导入
// @Tags	数据维护
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户编号"
// @Param	userName formData string true "用户昵称"
// @Param	password formData string false "用户密码"
// @Param	mobile formData string true "手机号"
// @Param	image formData string false "头像url"
// @Param	createTime formData string false "注册日期时间"
// @Param	status formData int true "用户状态 1-正常 -1 禁用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/user/transfer [post][get]
func UserTransfer(params map[string]string) mgresult.Result {
	status, _ := strconv.Atoi(params["status"])
	return service.NewUserService().UserTransfer(params["userId"], params["userName"], params["password"], params["mobile"], params["image"], params["createTime"], status)
}
