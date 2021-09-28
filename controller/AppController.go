package controller

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/service"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"strconv"
)

// AddApplication	godoc
// @Summary		添加新应用
// @Description	添加新应用
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appName formData string true "应用名称"
// @Param	descript formData string true "应用描述"
// @Param	smsSignCode formData string false "应用短信签名"
// @Param	smsLoginTemplate formData string false "应用登录短信模板，模板中短信验证码部分必须为{code}"
// @Param	smsRegisterTemplate formData string false "应用用户注册短信模板，模板中短信验证码部分必须为{code}"
// @Param	tokenTtl formData string true "本应用的token有效期，单位为分钟"
// @Param	maxOnline formData string true "本应用同时在线token数量"
// @Param	termTypes formData string true "本应用支持的终端类型，JSON数组。终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	verifySign formData int true "本应用是否验证签名，1-是 0-否"
// @Success 200 {string} string	"ok"
// @Router	/app/add [post][get]
func AddApplication(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appName") {
		return *mgresult.Error(-1, "应用名称不可为空")
	}
	if !utils.Exists(params, "descript") {
		return *mgresult.Error(-1, "应用简介不可为空")
	}
	if !utils.Exists(params, "smsSignCode") {
		params["smsSignCode"] = mgconfig.GetConfigString("goss.sms.signcode")
	}
	if !utils.Exists(params, "smsLoginTemplate") {
		params["smsLoginTemplate"] = mgconfig.GetConfigString("goss.sms.template.login")
	}
	if !utils.Exists(params, "smsRegisterTemplate") {
		params["smsRegisterTemplate"] = mgconfig.GetConfigString("goss.sms.template.register")
	}
	if !utils.Exists(params, "tokenTtl") {
		return *mgresult.Error(-1, "应用token的有效期不可为空")
	}
	tokenTtl, _ := strconv.Atoi(params["tokenTtl"])
	if !utils.Exists(params, "maxOnline") {
		return *mgresult.Error(-1, "应用同时最大在线数不可为空")
	}
	maxOnline, _ := strconv.Atoi(params["maxOnline"])
	if !utils.Exists(params, "termTypes") {
		return *mgresult.Error(-1, "应用支持的终端类型不可为空")
	}
	termTypes := make([]int, 0)
	utils.FromJSON(params["termTypes"], &termTypes)
	verifySign := 1
	if utils.Exists(params, "verifySign") {
		verifySign, _ = strconv.Atoi(params["verifySign"])
	}
	return service.NewAppService().AddApplication(params["appName"], params["descript"], params["smsSignCode"], params["smsLoginTemplate"], params["smsRegisterTemplate"], tokenTtl, maxOnline, termTypes, verifySign)
}

// GetAppInfo	godoc
// @Summary		查看应用详情
// @Description	查看应用详情
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Success 200 {string} string	"ok"
// @Router	/app/get [post][get]
func GetAppInfo(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewAppService().GetAppInfo(params["appId"])
}

// UpdateAppInfo	godoc
// @Summary		修改应用信息
// @Description	修改应用信息
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	appName formData string false "应用名称"
// @Param	descript formData string false "应用描述"
// @Success 200 {string} string	"ok"
// @Router	/app/update [post][get]
func UpdateAppInfo(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewAppService().UpdateAppInfo(params["appId"], params["appName"], params["descript"])
}

// ResetAppKey	godoc
// @Summary		重置应用密钥
// @Description	重置应用密钥
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	appKey formData string true "原应用密钥"
// @Success 200 {string} string	"ok"
// @Router	/app/key/reset [post][get]
func ResetAppKey(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewAppService().ResetAppKey(params["appId"], params["appKey"])
}

// GetAppSettings	godoc
// @Summary		查看应用常规设置
// @Description	查看应用常规设置
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Success 200 {string} string	"ok"
// @Router	/app/setting/get [post][get]
func GetAppSettings(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewAppService().GetAppSettings(params["appId"])
}

// UpdateAppSettings	godoc
// @Summary		修改应用常规设置
// @Description	修改应用常规设置
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	smsSignCode formData string false "应用短信签名"
// @Param	smsLoginTemplate formData string false "应用登录短信模板，模板中短信验证码部分必须为{code}"
// @Param	smsRegisterTemplate formData string false "应用用户注册短信模板，模板中短信验证码部分必须为{code}"
// @Param	tokenTtl formData int false "本应用令牌的最大有效时间，单位为秒"
// @Param	maxOnline formData int false "本应用同时在线最大令牌数"
// @Param	termTypes formData string false "本应用支持的终端类型列表，JSON数组格式,终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	verifySign formData int false "本应用是否验证签名，1-是 0-否"
// @Success 200 {string} string	"ok"
// @Router	/app/setting/update [post][get]
func UpdateAppSettings(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	tokenTtl := 0
	if utils.Exists(params, "tokenTtl") {
		tokenTtl, _ = strconv.Atoi(params["tokenTtl"])
	}
	maxOnline := 0
	if utils.Exists(params, "maxOnline") {
		maxOnline, _ = strconv.Atoi(params["maxOnline"])
	}
	termTypes := make([]int, 0)
	if utils.Exists(params, "termTypes") {
		utils.FromJSON(params["termTypes"], &termTypes)
	}
	verifySign := -1
	if utils.Exists(params, "verifySign") {
		verifySign, _ = strconv.Atoi(params["verifySign"])
	}
	return service.NewAppService().UpdateAppSettings(params["appId"], params["smsSignCode"], params["smsLoginTemplate"], params["smsRegisterTemplate"], tokenTtl, maxOnline, termTypes, verifySign)
}

// SetAppKickRule	godoc
// @Summary		设置应用互踢规则
// @Description	设置应用互踢规则
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	kickRule formData int true "互踢类型 1-同端互踢 2-完全互踢 3-允许多端登录，不互踢 4-指定端互踢 5-指定端不互踢 6-指定端分组互踢"
// @Param	termTypes formData string false "当kickRule为4或5时，互踢或不互踢的终端类型列表，JSON数组格式,终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	termTypeGroups formData string false "当kickRule为6时，互踢组列表，JSON数组"
// @Success 200 {string} string	"ok"
// @Router	/app/kickrule/set [post][get]
func SetAppKickRule(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "kickRule") {
		return *mgresult.Error(-1, "应用互踢类型不可为空")
	}
	kickRule, _ := strconv.Atoi(params["kickRule"])
	termTypes := make([]int, 0)
	if utils.Exists(params, "termTypes") {
		utils.FromJSON(params["termTypes"], &termTypes)
	}
	termTypeGroups := make([]string, 0)
	if utils.Exists(params, "termTypeGroups") {
		utils.FromJSON(params["termTypeGroups"], &termTypeGroups)
	}
	return service.NewAppService().SetAppKickRule(params["appId"], kickRule, termTypes, termTypeGroups)
}

// SetTermTypeGroup	godoc
// @Summary		设置应用互踢组
// @Description	设置应用互踢组
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	groupName formData string true "互踢组名称"
// @Param	termTypes formData string false "互踢组内的终端类型列表，JSON数组格式,终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Success 200 {string} string	"ok"
// @Router	/app/termtypegroup/add [post][get]
func SetTermTypeGroup(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "groupName") {
		return *mgresult.Error(-1, "互踢组名称不可为空")
	}
	termTypes := make([]int, 0)
	if utils.Exists(params, "termTypes") {
		utils.FromJSON(params["termTypes"], &termTypes)
	}
	return service.NewAppService().SetTermTypeGroup(params["appId"], params["groupName"], termTypes)
}

// DeleteTermTypeGroup	godoc
// @Summary		删除应用互踢组
// @Description	删除应用互踢组
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	termTypeGroup formData string true "互踢组编号"
// @Success 200 {string} string	"ok"
// @Router	/app/termtypegroup/del [post][get]
func DeleteTermTypeGroup(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	if !utils.Exists(params, "termTypeGroup") {
		return *mgresult.Error(-1, "互踢组名称不可为空")
	}
	return service.NewAppService().DeleteTermTypeGroup(params["appId"], params["termTypeGroup"])
}

// ListAppTermTypeGroups	godoc
// @Summary		列出应用所有的互踢组
// @Description	列出应用所有的互踢组
// @Tags	应用管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Success 200 {string} string	"ok"
// @Router	/app/termtypegroup/list [post][get]
func ListAppTermTypeGroups(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "appId") {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	return service.NewAppService().ListAppTermTypeGroups(params["appId"])
}
