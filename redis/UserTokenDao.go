package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"time"
)

//Token信息表
const (
	KEY_USER_TOKEN        = "user:token:"
	KEY_USER_TOKEN_EXPIRE = "user:expire:token:"
	KEY_USER_MYTOKEN      = "user:mytoken:"
)

func SaveToken(token model.TokenInfo, ttl time.Duration) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	token.CreateTime = utils.ToDateTimeString(time.Now())
	token.Status = constant.NORMAL
	err := goredis.Set(KEY_USER_TOKEN+token.Token, utils.ToJSON(token), 0).Err()
	if err != nil {
		return err
	}
	return goredis.Set(KEY_USER_TOKEN_EXPIRE+token.Token, ttl.Minutes(), ttl).Err()
}

func GetToken(token string) (model.TokenInfo, error) {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return model.TokenInfo{}, errors.New("Redis connection failed")
	}
	tokenInfoStr := goredis.Get(KEY_USER_TOKEN + token).Val()
	if tokenInfoStr == "" {
		return model.TokenInfo{}, errors.New("无此token")
	}
	var tokenInfo model.TokenInfo
	utils.FromJSON(tokenInfoStr, &tokenInfo)
	return tokenInfo, nil
}

//延期
func PostponeToken(token string, ttl time.Duration) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	return goredis.Expire(KEY_USER_TOKEN_EXPIRE+token, ttl).Err()
}

//刷新token
func RefleshToken(userId, token string) (model.TokenInfo, error) {
	tokenInfo, err := GetToken(token)
	if err != nil {
		return tokenInfo, err
	}
	if tokenInfo.UserId != userId || tokenInfo.Token == "" {
		return model.TokenInfo{}, errors.New("令牌错误或已失效")
	}
	tokenInfo.Token = utils.GetUUIDString()
	tokenInfo.CreateTime = utils.ToDateTimeString(time.Now())
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return tokenInfo, errors.New("Redis connection failed")
	}
	err = goredis.Set(KEY_USER_TOKEN+tokenInfo.Token, utils.ToJSON(token), 0).Err()
	if err != nil {
		return tokenInfo, err
	}
	ttl, _ := goredis.Get(KEY_USER_TOKEN_EXPIRE + token).Int64()
	goredis.Set(KEY_USER_TOKEN_EXPIRE+tokenInfo.Token, ttl, time.Duration(ttl)*time.Minute).Err()
	return tokenInfo, nil
}

func DeleteToken(token string) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	goredis.Del(KEY_USER_TOKEN_EXPIRE + token)
	return goredis.Del(KEY_USER_TOKEN + token).Err()
}

func GetUserTokens(userId string) (model.UserTokens, error) {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return model.UserTokens{}, errors.New("Redis connection failed")
	}
	myTokens := goredis.ZRevRangeByScore(KEY_USER_MYTOKEN+userId, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Val()
	var userTokens model.UserTokens
	userTokens.UserId = userId
	userTokens.Tokens = make([]model.TokenInfo, 0)
	for i := 0; i < len(myTokens); i++ {
		token := myTokens[i]
		tokenInfoStr := goredis.Get(KEY_USER_TOKEN + token).Val()
		if tokenInfoStr == "" {
			goredis.ZRem(KEY_USER_MYTOKEN+userId, token)
			continue
		}
		var tokenInfo model.TokenInfo
		utils.FromJSON(tokenInfoStr, &tokenInfo)
		userTokens.Tokens = append(userTokens.Tokens, tokenInfo)
	}
	return userTokens, nil
}

func AddUserTokens(userId, token string) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	member := &redis.Z{
		Member: token,
		Score:  float64(time.Now().UnixNano()),
	}
	return goredis.ZAdd(KEY_USER_MYTOKEN+userId, member).Err()
}

func RemoveOneFromUserTokens(userId, token string) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	return goredis.ZRem(KEY_USER_MYTOKEN+userId, token).Err()
}

func AddUserAppTokens(userId, appId, token string) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	return goredis.SAdd("user:app:token:"+userId+":"+appId, token).Err()
}

func RemoveOneFromUserAppTokens(userId, appId, token string) error {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	return goredis.SRem("user:app:token:"+userId+":"+appId, token).Err()
}

func ListTokensFromUserAppTokens(userId, appId string) []string {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return nil
	}
	return goredis.SMembers("user:app:token:" + userId + ":" + appId).Val()
}
