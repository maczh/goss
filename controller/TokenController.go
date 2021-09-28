package controller

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/service"
	"github.com/maczh/utils"
	"strconv"
)

// TokenAuthenticate	godoc
// @Summary		token验证是否在线并且有效
// @Description	token验证是否在线并且有效
// @Tags	令牌验证
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	token formData string true "校验的令牌"
// @Param	termType formData int true "终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/token/auth [post][get]
func TokenAuthenticate(params map[string]string) mgresult.Result {
	if !utils.Exists(params, "termType") {
		return *mgresult.Error(-1, "终端类型必传")
	}
	termType, _ := strconv.Atoi(params["termType"])
	return service.NewUserService().Authenticate(params["token"], params["appId"], termType)
}

// UserLogout	godoc
// @Summary		用户注销
// @Description	用户注销
// @Tags	注销/下线
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	token formData string true "用户令牌"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/logout [post][get]
func UserLogout(params map[string]string) mgresult.Result {
	return service.NewUserService().UserLogout(params["token"])
}

// SystemKickUser	godoc
// @Summary		系统强制下线
// @Description	系统强制下线
// @Tags	注销/下线
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户id"
// @Param	token formData string false "要强制下线的令牌，如果传此参数，则kickAppId和termType参数不生效"
// @Param	kickAppId formData string false "要强制下线的应用编码，若不传则踢除该用户所有令牌"
// @Param	termType formData int false "要踢的终端类型，若不传为所有终端类型 1-Web应用 2-Windows应用 3-macOS应用 4-ios终端 5-安卓终端 6-微信小程序 7-支付宝小程序 8-H5应用"
// @Param	invalidType formData int false "被踢的失效类型，可自定义，默认为3-系统强制下线"
// @Param	reason formData string false "被踢理由"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/quit [post][get]
func SystemKickUser(params map[string]string) mgresult.Result {
	termType := 0
	if utils.Exists(params, "termType") {
		termType, _ = strconv.Atoi(params["termType"])
	}
	invalidType := constant.IT_SYSTEM_KICKED
	if utils.Exists(params, "invalidType") {
		invalidType, _ = strconv.Atoi(params["invalidType"])
	}
	return service.NewUserService().SystemKickUser(params["appId"], params["userId"], params["token"], params["kickAppId"], termType, invalidType, params["reason"])
}

// ListUserTokensByApp	godoc
// @Summary		查看用户指定应用的所有令牌
// @Description	查看用户指定应用的所有令牌
// @Tags	令牌管理
// @Accept	x-www-form-urlencoded
// @Produce json
// @Param	appId formData string true "应用编号"
// @Param	userId formData string true "用户id"
// @Param	sign formData string true "签名"
// @Success 200 {string} string	"ok"
// @Router	/token/list [post][get]
func ListUserTokensByApp(params map[string]string) mgresult.Result {
	return service.NewUserService().ListUserTokensByApp(params["userId"], params["appId"])
}
