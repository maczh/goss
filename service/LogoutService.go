package service

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/logic"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/redis"
	"github.com/maczh/mgconfig"
	"strconv"
)

func (us *UserService) UserLogout(token string) mgresult.Result {
	if token == "" {
		return mgresult.Error(-1, "令牌不可为空")
	}
	err := logic.DeleteUserToken(token, "用户注销", constant.IT_USER_LOGOUT)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	return mgresult.Success(nil)
}

func (us *UserService) SystemKickUser(appId, userId, token, kickAppId string, termType, invalidType int, reason string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return mgresult.Error(-1, "应用编号错误")
	}
	if reason == "" {
		reason = "系统强制注销"
	}
	if token != "" {
		tokenInfo, err := redis.GetToken(token)
		if err != nil {
			return mgresult.Error(-1, err.Error())
		}
		if tokenInfo.Token == "" {
			return mgresult.Error(-1, "无此令牌")
		}
		if tokenInfo.AppId != appId {
			return mgresult.Error(-1, "非本应用令牌")
		}
		_ = logic.DeleteUserToken(token, reason, invalidType)
		return mgresult.Success(nil)
	}
	if kickAppId == "" {
		//踢除用户所有token
		tokens, err := redis.GetUserTokens(userId)
		if err != nil {
			return mgresult.Error(-1, err.Error())
		}
		for _, token := range tokens.Tokens {
			_ = logic.DeleteUserToken(token.Token, reason, invalidType)
		}
	} else {
		goredis := mgconfig.GetRedisConnection()
		defer mgconfig.ReturnRedisConnection(goredis)
		if goredis == nil {
			return mgresult.Error(-1, "数据库连接异常")
		}
		//中队用户指定应用内的token
		if termType > 0 {
			tokens := goredis.SMembers("user:app:" + appId + ":term:" + strconv.Itoa(termType) + ":token:" + userId).Val()
			for _, token := range tokens {
				_ = logic.DeleteUserToken(token, reason, invalidType)
			}
		} else {
			tokens := redis.ListTokensFromUserAppTokens(userId, appId)
			for _, token := range tokens {
				_ = logic.DeleteUserToken(token, reason, invalidType)
			}
		}
	}
	return mgresult.Success(nil)
}
