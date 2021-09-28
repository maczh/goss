package controller

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/service/thirdpart"
	"github.com/maczh/utils"
	"strconv"
)

// LoginThirdUser	godoc
// @Summary		用第三方认证登录
// @Description	用第三方认证登录
// @Tags	第三方登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	platform formData string true "第三方用户系统名称：wechat/qq/alipay/apple"
// @Param	thirdUserId formData string true "第三方用户唯一编号，微信的unionid,QQ的openid，支付宝的user_id"
// @Param	userIp formData string true "用户IP"
// @Param	deviceId formData string false "终端设备唯一代码"
// @Param	userAgent formData string true "终端客户端信息"
// @Param	deviceInfo formData string true "终端设备信息"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/login/oauth2 [post][get]
func LoginThirdUser(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "platform") || params["platform"] == "" {
		return *mgresult.Error(-1, "平台名称不可为空")
	}
	if !utils.Exists(params, "thirdUserId") || params["thirdUserId"] == "" {
		return *mgresult.Error(-1, "第三方用户编号不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "userIp") || params["userIp"] == "" {
		return *mgresult.Error(-1, "用户IP不可为空")
	}
	if !utils.Exists(params, "userAgent") || params["userAgent"] == "" {
		return *mgresult.Error(-1, "客户端信息不可为空")
	}
	if !utils.Exists(params, "deviceInfo") || params["deviceInfo"] == "" {
		return *mgresult.Error(-1, "终端设备信息不可为空")
	}
	if !utils.Exists(params, "termType") {
		return *mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return thirdpart.LoginThirdUser(params["platform"], params["thirdUserId"], params["appId"], params["deviceId"], params["userIp"], params["userAgent"], params["deviceInfo"], termType)
}

// BindThirdUser	godoc
// @Summary		绑定第三方认证用户
// @Description	绑定第三方认证用户
// @Tags	第三方登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	platform formData string true "第三方用户系统名称：wechat/qq/alipay/apple"
// @Param	thirdUserId formData string true "第三方用户唯一编号，微信的unionid,QQ的openid，支付宝的user_id"
// @Param	userId formData string false "系统用户编号"
// @Param	mobile formData string false "用户手机号码"
// @Param	nickName formData string false "第三方的用户昵称"
// @Param	sex formData string false "第三方的用户性别"
// @Param	province formData string false "第三方的用户省份"
// @Param	city formData string false "第三方的用户城市"
// @Param	country formData string false "第三方的用户国家"
// @Param	headImageUrl formData string false "第三方的用户头像图片url"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/bind/oauth2 [post][get]
func BindThirdUser(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "platform") || params["platform"] == "" {
		return *mgresult.Error(-1, "平台名称不可为空")
	}
	if !utils.Exists(params, "thirdUserId") || params["thirdUserId"] == "" {
		return *mgresult.Error(-1, "第三方用户编号不可为空")
	}
	return thirdpart.BindThirdUser(params["platform"], params["thirdUserId"], params["userId"], params["mobile"], params["nickName"], params["sex"], params["province"], params["city"], params["country"], params["headImageUrl"])
}

// UnBindThirdUser	godoc
// @Summary		第三方认证用户解除绑定
// @Description	第三方认证用户解除绑定
// @Tags	第三方登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	platform formData string true "第三方用户系统名称：wechat/qq/alipay/apple"
// @Param	userId formData string true "系统用户编号"
// @Param	thirdUserId formData string true "第三方用户唯一编号，微信的unionid,QQ的openid，支付宝的user_id"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/unbind/oauth2 [post][get]
func UnBindThirdUser(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "platform") || params["platform"] == "" {
		return *mgresult.Error(-1, "平台名称不可为空")
	}
	if !utils.Exists(params, "thirdUserId") || params["thirdUserId"] == "" {
		return *mgresult.Error(-1, "第三方用户编号不可为空")
	}
	if !utils.Exists(params, "userId") || params["userId"] == "" {
		return *mgresult.Error(-1, "用户编号不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return thirdpart.UnBindThirdUser(params["platform"], params["appId"], params["userId"], params["thirdUserId"])
}

// GetThirdUserInfo	godoc
// @Summary		获取第三方用户信息
// @Description	获取第三方用户信息
// @Tags	第三方登录
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	platform formData string true "第三方用户系统名称：wechat/qq/alipay/apple"
// @Param	userId formData string true "系统用户编号"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/user/oauth2 [post][get]
func GetThirdUserInfo(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "platform") || params["platform"] == "" {
		return *mgresult.Error(-1, "平台名称不可为空")
	}
	if !utils.Exists(params, "userId") || params["userId"] == "" {
		return *mgresult.Error(-1, "用户编号不可为空")
	}
	if !utils.Exists(params, "appId") || params["appId"] == "" {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return thirdpart.GetThirdUserInfo(params["platform"], params["userId"], params["appId"])
}
