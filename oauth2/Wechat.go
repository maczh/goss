package oauth2

import (
	"gitee.com/zchunshan/d3auth"
	"github.com/maczh/mgconfig"
)

type OAuth interface {
	SetSecret(appId, secret, redirectUrl string) *OAuth
	GetOAuthRedirectURL() string
	GetOAuthUserInfo(code string) (string, error)
}

type WechatPublic struct {
	WxConfig d3auth.Auth_conf
	WxAuth   *d3auth.Auth_wx
}

func (wxp *WechatPublic) SetSecret(appId, secret, redirectUrl string) *WechatPublic {
	if appId == "" {
		appId = mgconfig.GetConfigString("goss.oauth2.wechat.public.appid")
	}
	if secret == "" {
		secret = mgconfig.GetConfigString("goss.oauth2.wechat.public.secret")
	}
	wxp.WxConfig = d3auth.Auth_conf{
		Appid:  appId,
		Appkey: secret,
		Rurl:   redirectUrl,
	}
	return wxp
}

func (wxp *WechatPublic) GetOAuthRedirectURL() string {
	wxp.WxAuth = d3auth.NewAuth_wx(&wxp.WxConfig)
	return wxp.WxAuth.Get_Rurl("state")
}

func (wxp *WechatPublic) GetOAuthUserInfo(code string) (string, error) {
	wxres, err := wxp.WxAuth.Get_Token("code")
	if err != nil {
		return "", err
	}
	userInfo, err := wxp.WxAuth.Get_User_Info(wxres.Access_Token, wxres.Openid)
	return userInfo, err
}
