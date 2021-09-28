package constant

//登录方式类型
const (
	_                 = iota
	LT_TOKEN          //1
	LT_USERPWD        //2
	LT_SMS            //3
	LT_FIGNERPRINT    //4
	LT_FACEID         //5
	LT_WECHAT         //6
	LT_WECHAT_MINIAPP //7
	LT_ALIPAY         //8
	LT_QQ             //9
	LT_REGISTER_SMS   //10
	LT_REGISTER_OTHER //11
	LT_MOBILE_ONEKEY  //12
)
