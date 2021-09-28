package main

import (
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/maczh/gintool"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/aop"
	"github.com/maczh/goss/controller"
	_ "github.com/maczh/goss/docs"
	"github.com/maczh/mgtrace"
	"github.com/maczh/utils"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

/**
统一路由映射入口
*/
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	engine := gin.Default()

	//添加跟踪日志
	engine.Use(mgtrace.TraceId())

	//设置接口日志
	engine.Use(gintool.SetRequestLogger())
	//添加跨域处理
	engine.Use(gintool.Cors())

	//添加签名验证
	engine.Use(aop.VerifySign())

	//添加swagger支持
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//处理全局异常
	engine.Use(nice.Recovery(recoveryHandler))

	//设置404返回的内容
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, *mgresult.Error(-1, "请求的方法不存在"))
	})

	var result mgresult.Result
	//添加所需的路由映射
	//认证
	engine.Any("/token/auth", func(c *gin.Context) {
		result = controller.TokenAuthenticate(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//注册
	engine.Any("/register/sms", func(c *gin.Context) {
		result = controller.SendRegisterSms(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/register/sms/confirm", func(c *gin.Context) {
		result = controller.RegisterBySmsCode(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/register/pwd", func(c *gin.Context) {
		result = controller.RegisterByPassword(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//登录
	engine.Any("/login/sms", func(c *gin.Context) {
		result = controller.SendLoginSms(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/sms/confirm", func(c *gin.Context) {
		result = controller.LoginBySmsCode(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/pwd", func(c *gin.Context) {
		result = controller.LoginByPassword(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/finger/auth", func(c *gin.Context) {
		result = controller.GetFingerPrintCode(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/finger", func(c *gin.Context) {
		result = controller.LoginByFingerPrintCode(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/faceid/auth", func(c *gin.Context) {
		result = controller.GetFaceIdCode(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/faceid", func(c *gin.Context) {
		result = controller.LoginByFaceIdCode(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/login/mobile", func(c *gin.Context) {
		result = controller.LoginByAliMobile(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//注销/下线
	engine.Any("/logout", func(c *gin.Context) {
		result = controller.UserLogout(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/quit", func(c *gin.Context) {
		result = controller.SystemKickUser(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//应用管理
	engine.Any("/app/add", func(c *gin.Context) {
		result = controller.AddApplication(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/get", func(c *gin.Context) {
		result = controller.GetAppInfo(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/update", func(c *gin.Context) {
		result = controller.UpdateAppInfo(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/key/reset", func(c *gin.Context) {
		result = controller.ResetAppKey(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/setting/get", func(c *gin.Context) {
		result = controller.GetAppSettings(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/setting/update", func(c *gin.Context) {
		result = controller.UpdateAppSettings(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/kickrule/set", func(c *gin.Context) {
		result = controller.SetAppKickRule(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/termtypegroup/add", func(c *gin.Context) {
		result = controller.SetTermTypeGroup(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/termtypegroup/del", func(c *gin.Context) {
		result = controller.DeleteTermTypeGroup(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/app/termtypegroup/list", func(c *gin.Context) {
		result = controller.ListAppTermTypeGroups(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//用户管理
	engine.Any("/user/get", func(c *gin.Context) {
		result = controller.GetuserInfo(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/user/update", func(c *gin.Context) {
		result = controller.UpdateUserInfo(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/user/status/update", func(c *gin.Context) {
		result = controller.UpdateUserStatus(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/user/transfer", func(c *gin.Context) {
		result = controller.UserTransfer(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//令牌管理
	engine.Any("/token/list", func(c *gin.Context) {
		result = controller.ListUserTokensByApp(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	//第三方
	engine.Any("/login/oauth2", func(c *gin.Context) {
		result = controller.LoginThirdUser(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/bind/oauth2", func(c *gin.Context) {
		result = controller.BindThirdUser(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/unbind/oauth2", func(c *gin.Context) {
		result = controller.UnBindThirdUser(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	engine.Any("/user/oauth2", func(c *gin.Context) {
		result = controller.GetThirdUserInfo(utils.GinParamMap(c))
		c.JSON(http.StatusOK, result)
	})

	return engine
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusOK, *mgresult.Error(-1, "系统异常，请联系客服"))
}
