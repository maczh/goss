package logic

import (
	"errors"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/redis"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgtrace"
	"github.com/maczh/utils"
	"strconv"
	"time"
)

func KickToken(token, byToken, reason string) error {
	kickedToken, err := redis.GetToken(token)
	if err != nil {
		return err
	}
	newToken := model.TokenInfo{UserId: "system", Token: "system"}
	if byToken != "" {
		newToken, err = redis.GetToken(byToken)
		if err != nil {
			return err
		}
	}
	//删除Redis中被踢的token信息
	err = redis.DeleteToken(token)
	if err != nil {
		return err
	}
	kickedToken.Status = constant.IT_KICKED
	kickedToken.Error = "被踢下线"
	if byToken != "" {
		kickedToken.Error = "您的账号在其他终端登录，您必须重新登录，若非本人操作请尽快修改密码或联系客服"
	} else {
		kickedToken.Error = reason
	}
	//保存被踢日志
	userKickLog := model.UserKickLog{
		Token:       token,
		RequestId:   mgtrace.GetRequestId(),
		KickedToken: kickedToken,
		ByNewToken:  newToken,
	}
	userKickLog, err = mongo.InsertKickLog(userKickLog)
	//删除用户所有Token表
	redis.RemoveOneFromUserTokens(kickedToken.UserId, token)
	//记录失效日志
	invalidInfo := model.TokenInvalidInfo{
		InvalidType: constant.IT_KICKED,
		Message:     reason,
	}
	invalidLog := model.TokenInvalidLog{
		Token:       token,
		TokenInfo:   kickedToken,
		CreateTime:  kickedToken.CreateTime,
		InvalidTime: utils.ToDateTimeString(time.Now()),
		InvalidInfo: invalidInfo,
	}
	_, err = mongo.InsertInvalidLog(invalidLog)
	return err
}

func AddNewTokenWithKickRules(tokenInfo model.TokenInfo) error {
	//提取互踢规则
	appKickRule, err := redis.GetAppKickedRule(tokenInfo.AppId)
	if err != nil {
		return err
	}
	kickedTokenList := make([]string, 0)
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	appSettings, err := mongo.GetAppSettings(tokenInfo.AppId)
	if err != nil {
		logs.Error("获取应用设置错误:{}", err.Error())
		return err
	}
	if !utils.SliceContainsInt(appSettings.TermTypes, tokenInfo.TermType) {
		logs.Error("此应用不支持此终端类型，如需添加请联系客服人员")
		return errors.New("此应用不支持此终端类型")
	}
	switch appKickRule.KickRule {
	//同端互踢
	case constant.KICK_SAME_TERMINAL:
		tokens := goredis.SMembers("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(tokenInfo.TermType) + ":token:" + tokenInfo.UserId).Val()
		kickedTokenList = append(kickedTokenList, tokens...)
		goredis.Del("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(tokenInfo.TermType) + ":token:" + tokenInfo.UserId)
		goredis.SAdd("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, tokenInfo.Token)
		for _, tk := range tokens {
			redis.RemoveOneFromUserTokens(tokenInfo.UserId, tk)
		}
	//完全互踢
	case constant.ONLY_ONE_TERMINAL:
		allTokens, err := redis.GetUserTokens(tokenInfo.UserId)
		if err != nil {
			logs.Error("获取用户所有token时异常:{}", err.Error())
			return err
		}
		for _, tkInfo := range allTokens.Tokens {
			kickedTokenList = append(kickedTokenList, tkInfo.Token)
			goredis.Del("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(tkInfo.TermType) + ":token:" + tokenInfo.UserId)
			redis.RemoveOneFromUserTokens(tokenInfo.UserId, tkInfo.Token)
		}
		goredis.SAdd("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, tokenInfo.Token)
	//不互踢,仅限制终端数，先进先出
	case constant.ALLOW_MULTIPLE_ONLINE:
		goredis.SAdd("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, tokenInfo.Token)
		if appSettings.MaxOnlineTokens > 0 {
			allTokens, err := redis.GetUserTokens(tokenInfo.UserId)
			if err != nil {
				logs.Error("获取用户所有token时异常:{}", err.Error())
				return err
			}
			if len(allTokens.Tokens) >= appSettings.MaxOnlineTokens {
				for _, token := range allTokens.Tokens[appSettings.MaxOnlineTokens-2:] {
					kickedTokenList = append(kickedTokenList, token.Token)
					goredis.Del("user:app:" + token.AppId + ":term:" + strconv.Itoa(token.TermType) + ":token:" + token.UserId)
					redis.RemoveOneFromUserTokens(tokenInfo.UserId, token.Token)
				}
			}
		}
	//指定终端类型之间互踢
	case constant.ONLY_ONE_WITHIN_TERMTYPES:
		if utils.SliceContainsInt(appKickRule.TermTypes, tokenInfo.TermType) {
			for _, termType := range appKickRule.TermTypes {
				tokens := goredis.SMembers("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(termType) + ":token:" + tokenInfo.UserId).Val()
				if len(tokens) > 0 {
					for _, token := range tokens {
						kickedTokenList = append(kickedTokenList, token)
						goredis.Del("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(termType) + ":token:" + tokenInfo.UserId)
						redis.RemoveOneFromUserTokens(tokenInfo.UserId, token)
					}
				}
			}
		}
		goredis.SAdd("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, tokenInfo.Token)
	//指定终端类型之外互踢，之间不互踢
	case constant.ONLY_ONE_WITHOUT_TERMTYPES:
		if !utils.SliceContainsInt(appKickRule.TermTypes, tokenInfo.TermType) {
			for _, termType := range appSettings.TermTypes {
				if !utils.SliceContainsInt(appKickRule.TermTypes, termType) {
					tokens := goredis.SMembers("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(termType) + ":token:" + tokenInfo.UserId).Val()
					if len(tokens) > 0 {
						for _, token := range tokens {
							kickedTokenList = append(kickedTokenList, token)
							goredis.Del("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(termType) + ":token:" + tokenInfo.UserId)
							redis.RemoveOneFromUserTokens(tokenInfo.UserId, token)
						}
					}
				}
			}
		}
		goredis.SAdd("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, tokenInfo.Token)
	//指定终端类型分组内互踢
	case constant.ONLY_ONE_WITHIN_TERMTYPE_GROUPS:
		for _, termTypeGroup := range appKickRule.TermTypeGroups {
			if utils.SliceContainsInt(termTypeGroup.TermTypes, tokenInfo.TermType) {
				for _, termType := range termTypeGroup.TermTypes {
					tokens := goredis.SMembers("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(termType) + ":token:" + tokenInfo.UserId).Val()
					if len(tokens) > 0 {
						for _, token := range tokens {
							kickedTokenList = append(kickedTokenList, token)
							goredis.Del("user:app:" + tokenInfo.AppId + ":term:" + strconv.Itoa(termType) + ":token:" + tokenInfo.UserId)
							redis.RemoveOneFromUserTokens(tokenInfo.UserId, token)
						}
					}
				}
			}
		}
		goredis.SAdd("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, tokenInfo.Token)
	}
	for _, token := range kickedTokenList {
		err := KickToken(token, tokenInfo.Token, "终端登录互踢")
		if err != nil {
			logs.Error("踢掉{}错误:{}", token, err.Error())
		}
	}
	return err
}
