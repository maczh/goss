package aop

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gin-gonic/gin"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/logs"
	"github.com/maczh/utils"
	"strings"
)

func VerifySign() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifySign(c)
	}
}

func verifySign(ctx *gin.Context) {
	if strings.Contains(ctx.Request.RequestURI, "/app/") || strings.Contains(ctx.Request.RequestURI, "/docs/") {
		return
	}
	params := utils.GinParamMap(ctx)
	if !utils.Exists(params, "appId") {
		ctx.AbortWithStatusJSON(200, mgresult.Error(-1, "应用编号不可为空"))
		return
	}
	appSettings, err := mongo.GetAppSettings(params["appId"])
	if err != nil {
		ctx.AbortWithStatusJSON(200, mgresult.Error(-1, err.Error()))
		return
	}
	if appSettings.AppId == "" {
		ctx.AbortWithStatusJSON(200, mgresult.Error(-1, "应用代码不正确"))
		return
	}
	if appSettings.VerifySign == false {
		return
	}
	if !utils.Exists(params, "sign") {
		ctx.AbortWithStatusJSON(200, mgresult.Error(-1, "应用签名不可为空"))
		return
	}
	appInfo, err := mysql.GetAppInfoByAppId(params["appId"])
	if err != nil {
		ctx.AbortWithStatusJSON(200, mgresult.Error(-1, err.Error()))
		return
	}
	sortmap := treemap.NewWithStringComparator()
	for k, v := range params {
		if v != "" && k != "sign" {
			sortmap.Put(k, v)
		}
	}
	querystr := ""
	sortmap.Each(func(key interface{}, value interface{}) {
		querystr = querystr + key.(string) + "=" + value.(string) + "&"
	})
	querystr = querystr + "appKey=" + appInfo.AppKey
	sign := utils.MD5Encode(querystr)
	logs.Debug("签名明文:{},结果:{}", querystr, sign)
	if params["sign"] != sign {
		ctx.AbortWithStatusJSON(200, mgresult.Error(-1, "应用签名不正确"))
	}
}
