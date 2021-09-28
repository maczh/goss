package service

import (
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/redis"
	"github.com/maczh/utils"
)

type FaceIdService struct{}

func NewFaceIdService() *FaceIdService {
	return &FaceIdService{}
}

func (s *FaceIdService) GetFaceIdCode(userId, token string) mgresult.Result {
	tokenInfo, err := redis.GetToken(token)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if tokenInfo.UserId != userId {
		return *mgresult.Error(-1, "令牌错误")
	}
	userSpecialCode, err := mongo.GetUserSpecialCode(userId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if userSpecialCode.FaceIdCode == "" {
		userSpecialCode.UserId = userId
		userSpecialCode.FaceIdCode = utils.GetRandomCaseString(64)
		userSpecialCode, _ = mongo.UpsertUserSpecialCode(userSpecialCode)
	}
	return *mgresult.Success(userSpecialCode)
}

func (s *FaceIdService) LoginByFaceIdCode(mobile, faceIdCode, appId, deviceId, userIp, userAgent, deviceInfo string, termType int) mgresult.Result {
	userInfo, err := mysql.GetUserInfoByMobile(mobile)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if userInfo.UserId == "" {
		return *mgresult.Error(-1, "该用户不存在")
	}
	if userInfo.Status != constant.NORMAL {
		return *mgresult.Error(-1, "该账户状态不可用")
	}
	b, err := NewAppService().CheckAppTermType(appId, termType)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if !b {
		return *mgresult.Error(-1, "本应用不支持此终端类型登录")
	}
	userSpecialCode, err := mongo.GetUserSpecialCode(userInfo.UserId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if userSpecialCode.FaceIdCode == "" {
		return *mgresult.Error(-1, "未开通脸部识别登录授权，请先开通脸部识别登录授权")
	}
	if userSpecialCode.FaceIdCode != faceIdCode {
		return *mgresult.Error(-1, "脸部识别码不匹配")
	}
	result, err := NewUserService().loadUser("", mobile).login(appId, deviceId, userIp, userAgent, deviceInfo, termType, constant.LT_FACEID)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	//返回tokenInfo
	return *mgresult.Success(result)
}
