package nacos

import (
	"errors"
	"github.com/maczh/gintool/mgresult"
	"github.com/maczh/logs"
	"github.com/maczh/mgcall"
	"github.com/maczh/utils"
)

const (
	SERVICE_PUBLIC_SERVICE = "publicservice"
	URI_SMS                = "/sms/send"
	URI_MOBILE_GET         = "/alm/mobile/get"
)

func SendSms(mobile, templatecode, signcode, json string) error {
	params := make(map[string]string)
	params["mobile"] = mobile
	params["templatecode"] = templatecode
	params["signcode"] = signcode
	params["json"] = json
	res, err := mgcall.Call(SERVICE_PUBLIC_SERVICE, URI_SMS, params)
	if err != nil {
		logs.Error("微服务{}{}调用异常:{}", SERVICE_PUBLIC_SERVICE, URI_SMS, err.Error())
		return err
	}
	var result mgresult.Result
	utils.FromJSON(res, &result)
	if result.Status == 1 {
		return nil
	} else {
		return errors.New(result.Msg)
	}
}

func GetMobileNumber(requestId, mobileToken string) (string, error) {
	params := make(map[string]string)
	params["requestId"] = requestId
	params["token"] = mobileToken
	res, err := mgcall.Call(SERVICE_PUBLIC_SERVICE, URI_MOBILE_GET, params)
	if err != nil {
		logs.Error("微服务{}{}调用异常:{}", SERVICE_PUBLIC_SERVICE, URI_MOBILE_GET, err.Error())
		return "", err
	}
	var result mgresult.Result
	utils.FromJSON(res, &result)
	if result.Status == 1 {
		data := result.Data.(map[string]string)
		return data["mobile"], nil
	} else {
		return "", errors.New(result.Msg)
	}
}
