package service

import (
	"errors"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/mongo"
	"github.com/maczh/goss/mysql"
	"github.com/maczh/goss/redis"
	"github.com/maczh/logs"
	"github.com/maczh/utils"
)

type AppService struct{}

func NewAppService() *AppService {
	return &AppService{}
}

func (app *AppService) CheckAppTermType(appId string, termType int) (bool, error) {
	appSettings, err := mongo.GetAppSettings(appId)
	if err != nil {
		logs.Error("获取应用设置错误:{}", err.Error())
		return false, err
	}
	if appSettings.AppId == "" {
		return false, errors.New("应用代码错误")
	}
	if utils.SliceContainsInt(appSettings.TermTypes, termType) {
		return true, nil
	} else {
		return false, nil
	}
}

func (app *AppService) AddApplication(appName, descript, smsSignCode, smsLoginTemplate, smsRegisterTemplate string, tokenTtl, maxOnline int, termTypes []int, verifySign int) mgresult.Result {
	var err error
	appInfo := model.AppInfo{
		AppName:  appName,
		Descript: descript,
	}
	appInfo, err = mysql.InsertAppInfo(appInfo)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if tokenTtl == 0 {
		//默认token有效期为10天
		tokenTtl = 10 * 24 * 60
	}
	if maxOnline == 0 {
		//默认最大10个token
		maxOnline = 10
	}
	if termTypes == nil || len(termTypes) == 0 {
		termTypes = []int{
			constant.WEB,
			constant.PC,
			constant.MAC,
			constant.H5,
			constant.ANDROID,
			constant.IOS,
			constant.WECHAT,
		}
	}
	appSettings := model.AppSettings{
		AppId:               appInfo.AppId,
		TokenTtl:            tokenTtl,
		TermTypes:           termTypes,
		MaxOnlineTokens:     maxOnline,
		SmsSignCode:         smsSignCode,
		SmsLoginTemplate:    smsLoginTemplate,
		SmsRegisterTemplate: smsRegisterTemplate,
		VerifySign:          verifySign == 1,
	}
	_, err = mongo.InsertAppSettings(appSettings)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(appInfo)
}

func (app *AppService) GetAppInfo(appId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	return *mgresult.Success(appInfo)
}

func (app *AppService) UpdateAppInfo(appId, appName, descript string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	if appName != "" {
		appInfo.AppName = appName
	}
	if descript != "" {
		appInfo.Descript = descript
	}
	err = mysql.UpdateAppInfo(appInfo)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(appInfo)
}

func (app *AppService) ResetAppKey(appId, appKey string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	if appInfo.AppKey != appKey {
		return *mgresult.Error(-1, "原应用密钥验证错误")
	}
	appInfo.AppKey = utils.GetRandomCaseString(32)
	err = mysql.UpdateAppInfo(appInfo)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(appInfo)
}

func (app *AppService) GetAppSettings(appId string) mgresult.Result {
	appSettings, err := mongo.GetAppSettings(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appSettings.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	return *mgresult.Success(appSettings)
}

func (app *AppService) UpdateAppSettings(appId, smsSignCode, smsLoginTemplate, smsRegisterTemplate string, tokenTtl, maxOnline int, termTypes []int, verifySign int) mgresult.Result {
	appSettings, err := mongo.GetAppSettings(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appSettings.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	if tokenTtl != 0 {
		appSettings.TokenTtl = tokenTtl
	}
	if maxOnline != 0 {
		appSettings.MaxOnlineTokens = maxOnline
	}
	if termTypes != nil && len(termTypes) > 0 {
		appSettings.TermTypes = termTypes
	}
	if smsSignCode != "" {
		appSettings.SmsSignCode = smsSignCode
	}
	if smsLoginTemplate != "" {
		appSettings.SmsLoginTemplate = smsLoginTemplate
	}
	if smsRegisterTemplate != "" {
		appSettings.SmsRegisterTemplate = smsRegisterTemplate
	}
	if verifySign != -1 {
		appSettings.VerifySign = verifySign == 1
	}
	err = mongo.UpdateAppSettings(appSettings)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(appSettings)
}

func (app *AppService) SetAppKickRule(appId string, kickRule int, termTypes []int, termTypeGroups []string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	if kickRule <= 0 || kickRule > constant.ONLY_ONE_WITHIN_TERMTYPE_GROUPS {
		return *mgresult.Error(-1, "互踢类型错误")
	}
	if (kickRule == constant.ONLY_ONE_WITHIN_TERMTYPES || kickRule == constant.ONLY_ONE_WITHOUT_TERMTYPES) && (termTypes == nil || len(termTypes) == 0) {
		return *mgresult.Error(-1, "互踢终端类型列表不可为空")
	}
	if kickRule == constant.ONLY_ONE_WITHIN_TERMTYPE_GROUPS {
		if termTypeGroups == nil || len(termTypeGroups) == 0 {
			return *mgresult.Error(-1, "互踢组列表不可为空")
		}
		if !redis.IsAppTermTypeGroup(appId, termTypeGroups) {
			return *mgresult.Error(-1, "互踢组的编号错误")
		}
	}
	appKickRule := model.AppKickRule{
		AppId:          appId,
		KickRule:       kickRule,
		TermTypes:      termTypes,
		TermTypeGroups: redis.GetTermTypeGroups(termTypeGroups),
	}
	err = redis.SetAppKickRule(appKickRule)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(appKickRule)
}

func (app *AppService) SetTermTypeGroup(appId, groupName string, termTypes []int) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	if termTypes == nil || len(termTypes) == 0 {
		return *mgresult.Error(-1, "互踢组终端类型列表不可为空")
	}
	if groupName == "" {
		return *mgresult.Error(-1, "互踢组名称不可为空")
	}
	termTypeGroup := model.TermTypeGroup{
		GroupName: groupName,
		TermTypes: termTypes,
	}
	termTypeGroup, err = redis.SetTermTypeGroup(appId, termTypeGroup)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(termTypeGroup)
}

func (app *AppService) DeleteTermTypeGroup(appId, termTypeGroup string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	if termTypeGroup == "" {
		return *mgresult.Error(-1, "互踢组ID不可为空")
	}
	err = redis.DelTermTypeGroup(appId, termTypeGroup)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(nil)
}

func (app *AppService) ListAppTermTypeGroups(appId string) mgresult.Result {
	appInfo, err := mysql.GetAppInfoByAppId(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	if appInfo.AppId == "" {
		return *mgresult.Error(-1, "应用编号错误")
	}
	termTypeGroups, err := redis.ListAppTermTypeGroups(appId)
	if err != nil {
		return *mgresult.Error(-1, err.Error())
	}
	return *mgresult.Success(termTypeGroups)
}
