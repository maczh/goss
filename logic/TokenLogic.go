package logic

import (
	"errors"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/redis"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"strconv"
	"time"
)

func NewUserToken(userId, appId, deviceId, userAgent, userIp, deviceInfo string, termType int) (model.TokenInfo, error) {
	//判断是否应用支持的终端类型
	appSettings, err := mongo.GetAppSettings(appId)
	if err != nil {
		return model.TokenInfo{}, err
	}
	if !utils.SliceContainsInt(appSettings.TermTypes, termType) {
		return model.TokenInfo{}, errors.New("本应用不支持此终端类型")
	}
	//若无设备Id，则生成一个
	if deviceId == "" {
		deviceId = utils.MD5Encode(userId + deviceInfo + userAgent)
	}
	//若用户设备库中无此设备，则自动添加
	userDevice, err := mongo.GetUserDevice(userId, deviceId)
	if err != nil {
		return model.TokenInfo{}, err
	}
	if userDevice.DeviceId == "" {
		userDevice.UserId = userId
		userDevice.DeviceId = deviceId
		userDevice.DeviceInfo = deviceInfo
		userDevice, err = mongo.InsertUserDevice(userDevice)
		if err != nil {
			return model.TokenInfo{}, err
		}
	}
	//生成token并保存
	tokenInfo := model.TokenInfo{
		Token:      utils.GetUUIDString(),
		UserId:     userId,
		AppId:      appId,
		DeviceId:   deviceId,
		UserAgent:  userAgent,
		UserIp:     userIp,
		TermType:   termType,
		Status:     constant.NORMAL,
		CreateTime: utils.ToDateTimeString(time.Now()),
	}
	err = redis.SaveToken(tokenInfo, time.Duration(appSettings.TokenTtl)*time.Minute)
	if err != nil {
		return model.TokenInfo{}, err
	}
	//保存用户应用终端token表,按互踢规则互踢
	err = AddNewTokenWithKickRules(tokenInfo)
	if err != nil {
		return model.TokenInfo{}, err
	}
	//保存用户所有Token表
	err = redis.AddUserTokens(userId, tokenInfo.Token)
	_ = redis.AddUserAppTokens(userId, appId, tokenInfo.Token)
	if err != nil {
		return model.TokenInfo{}, err
	}
	return tokenInfo, nil
}

func DeleteUserToken(token, reason string, invalidType int) error {
	//先取出原token内容
	tokenInfo, err := redis.GetToken(token)
	if err != nil {
		return err
	}
	if tokenInfo.Token == "" {
		return errors.New("此令牌不存在或已失效")
	}
	//删除token表
	redis.DeleteToken(token)
	//删除应用token表
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return errors.New("Redis connection failed")
	}
	goredis.SRem("user:app:"+tokenInfo.AppId+":term:"+strconv.Itoa(tokenInfo.TermType)+":token:"+tokenInfo.UserId, token)
	//删除用户所有Token表
	redis.RemoveOneFromUserTokens(tokenInfo.UserId, token)
	redis.RemoveOneFromUserAppTokens(tokenInfo.UserId, tokenInfo.AppId, token)
	//记录失效日志
	tokenInfo.Status = constant.IT_EXPIRED
	tokenInfo.Error = reason
	invalidInfo := model.TokenInvalidInfo{
		InvalidType: invalidType,
		Message:     reason,
	}
	invalidLog := model.TokenInvalidLog{
		Token:       token,
		TokenInfo:   tokenInfo,
		CreateTime:  tokenInfo.CreateTime,
		InvalidTime: utils.ToDateTimeString(time.Now()),
		InvalidInfo: invalidInfo,
	}
	_, err = mongo.InsertInvalidLog(invalidLog)
	return err
}

func VerifyToken(token, appId string, termType int) (string, bool, error, model.TokenInvalidLog) {
	tokenInfo, _ := redis.GetToken(token)
	//if err != nil {
	//	return false, err, model.TokenInvalidLog{}
	//}
	if tokenInfo.Token == "" {
		tokenInvalidLog, err := mongo.GetTokenInvalidLog(token)
		if err != nil {
			return "", false, err, tokenInvalidLog
		}
		if tokenInvalidLog.Token == "" {
			//无失效记录，应为不存在的token
			return "", false, errors.New("错误的令牌"), model.TokenInvalidLog{}
		} else {
			return "", false, nil, tokenInvalidLog
		}
	}
	if tokenInfo.Status == constant.NORMAL && tokenInfo.AppId == appId && tokenInfo.TermType == termType {
		return tokenInfo.UserId, true, nil, model.TokenInvalidLog{}
	}
	return "", false, errors.New("令牌不匹配"), model.TokenInvalidLog{}
}
