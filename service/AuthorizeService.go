package service

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/logic"
	"github.com/maczh/goss/model/response"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/redis"
	"github.com/maczh/logs"
	"time"
)

func (us *UserService) Authenticate(token, appId string, termType int) mgresult.Result {
	if token == "" {
		return *mgresult.Error(-1, "令牌不可为空")
	}
	if appId == "" {
		return *mgresult.Error(-1, "应用编号不可为空")
	}
	if termType == 0 {
		return *mgresult.Error(-1, "终端类型错误")
	}
	userId, online, err, invalidLog := logic.VerifyToken(token, appId, termType)
	validResult := response.TokenVerify{
		Token:       token,
		UserId:      userId,
		Status:      1,
		InvalidType: 0,
		Message:     "",
	}
	if online {
		appSettings, err := mongo.GetAppSettings(appId)
		if err != nil {
			logs.Error("获取应用设置错误:{}", err.Error())
		}
		if appSettings.AppId == appId {
			redis.PostponeToken(token, time.Duration(appSettings.TokenTtl)*time.Minute)
		}
		return *mgresult.Success(validResult)
	}
	if err != nil {
		validResult.Status = 2
		validResult.Message = err.Error()
		return *mgresult.Success(validResult)
	}
	if invalidLog.Token != "" {
		validResult.Status = 3
		validResult.Message = invalidLog.InvalidInfo.Message
		validResult.InvalidType = invalidLog.InvalidInfo.InvalidType
		return *mgresult.Success(validResult)
	}
	return *mgresult.Error(-1, err.Error())
}
