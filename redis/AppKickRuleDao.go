package redis

import (
	"errors"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"strconv"
)

const (
	KEY_KICK_RULE           = "app:kickrule:type:"     //互踢规则类型 1-同端互踢 2-完全互踢，单端在线 3-指定端互踢 4-指定端不互踢
	KEY_KICK_RULE_TERMTYPE  = "app:kickrule:termtype:" //指定端互踢或不互踢的终端类型清单
	KEY_KICK_TERMTYPE_GROUP = "app:kickrule:group:"    //互踢组
	KEY_APP_TERM_TYPE_GROUP = "app:termtype:group:"    //应用包含的互踢组集合
)

func GetAppKickedRule(appId string) (model.AppKickRule, error) {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return model.AppKickRule{}, errors.New("Redis connection failed")
	}
	var appKickRule model.AppKickRule
	appKickRule.AppId = appId
	appKickRule.KickRule, _ = goredis.Get(KEY_KICK_RULE + appId).Int()
	switch appKickRule.KickRule {
	case constant.ONLY_ONE_WITHIN_TERMTYPES,
		constant.ONLY_ONE_WITHOUT_TERMTYPES:
		termTypes := goredis.SMembers(KEY_KICK_RULE_TERMTYPE + appId).Val()
		appKickRule.TermTypes = make([]int, len(termTypes))
		for i := 0; i < len(termTypes); i++ {
			appKickRule.TermTypes[i], _ = strconv.Atoi(termTypes[i])
		}
	case constant.ONLY_ONE_WITHIN_TERMTYPE_GROUPS:
		termTypeGroups := goredis.SMembers(KEY_KICK_RULE_TERMTYPE + appId).Val()
		appKickRule.TermTypeGroups = make([]model.TermTypeGroup, len(termTypeGroups))
		for i := 0; i < len(termTypeGroups); i++ {
			termTypeGroup := goredis.Get(KEY_KICK_TERMTYPE_GROUP + termTypeGroups[i]).Val()
			utils.FromJSON(termTypeGroup, &appKickRule.TermTypeGroups[i])
		}
	default:
	}
	return appKickRule, nil
}

func SetAppKickRule(appKickRule model.AppKickRule) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	goredis.Set(KEY_KICK_RULE+appKickRule.AppId, appKickRule.KickRule, 0)
	switch appKickRule.KickRule {
	case constant.ONLY_ONE_WITHIN_TERMTYPES,
		constant.ONLY_ONE_WITHOUT_TERMTYPES:
		goredis.Del(KEY_KICK_RULE_TERMTYPE + appKickRule.AppId)
		goredis.SAdd(KEY_KICK_RULE_TERMTYPE+appKickRule.AppId, appKickRule.TermTypes)
	case constant.ONLY_ONE_WITHIN_TERMTYPE_GROUPS:
		goredis.Del(KEY_KICK_RULE_TERMTYPE + appKickRule.AppId)
		termTypeGroups := make([]string, len(appKickRule.TermTypeGroups))
		for i, group := range appKickRule.TermTypeGroups {
			termTypeGroups[i] = group.GroupId
		}
		goredis.SAdd(KEY_KICK_RULE_TERMTYPE+appKickRule.AppId, termTypeGroups)
	}
	return nil
}

func SetTermTypeGroup(appId string, termTypeGroup model.TermTypeGroup) (model.TermTypeGroup, error) {
	if termTypeGroup.GroupId == "" {
		termTypeGroup.GroupId = utils.GetRandomHexString(16)
	}
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return model.TermTypeGroup{}, errors.New("Redis connection failed")
	}
	goredis.Set(KEY_KICK_TERMTYPE_GROUP+termTypeGroup.GroupId, utils.ToJSON(termTypeGroup), 0)
	goredis.SAdd(KEY_APP_TERM_TYPE_GROUP+appId, termTypeGroup.GroupId)
	return termTypeGroup, nil
}

func DelTermTypeGroup(appId, termTypeGroupId string) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	goredis.Del(KEY_KICK_TERMTYPE_GROUP + termTypeGroupId)
	goredis.SRem(KEY_APP_TERM_TYPE_GROUP+appId, termTypeGroupId)
	return nil
}

func IsAppTermTypeGroup(appId string, termTypeGroups []string) bool {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return false
	}
	for _, termTypeGroup := range termTypeGroups {
		if !goredis.SIsMember(KEY_APP_TERM_TYPE_GROUP+appId, termTypeGroup).Val() {
			return false
		}
	}
	return true
}

func GetTermTypeGroups(termTypeGroups []string) []model.TermTypeGroup {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return nil
	}
	appTermTypeGroups := make([]model.TermTypeGroup, len(termTypeGroups))
	for i := 0; i < len(termTypeGroups); i++ {
		termTypeGroup := goredis.Get(KEY_KICK_TERMTYPE_GROUP + termTypeGroups[i]).Val()
		utils.FromJSON(termTypeGroup, &appTermTypeGroups[i])
	}
	return appTermTypeGroups
}

func ListAppTermTypeGroups(appId string) ([]model.TermTypeGroup, error) {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return nil, errors.New("Redis connection failed")
	}
	typeGroups := goredis.SMembers(KEY_APP_TERM_TYPE_GROUP + appId).Val()
	appTermTypeGroups := make([]model.TermTypeGroup, len(typeGroups))
	for i := 0; i < len(typeGroups); i++ {
		termTypeGroup := goredis.Get(KEY_KICK_TERMTYPE_GROUP + typeGroups[i]).Val()
		utils.FromJSON(termTypeGroup, &appTermTypeGroups[i])
	}
	return appTermTypeGroups, nil
}
