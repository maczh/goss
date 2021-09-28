package constant

const (
	_                               = iota //互踢规则类型
	KICK_SAME_TERMINAL                     //1-同端互踢
	ONLY_ONE_TERMINAL                      //2-完全互踢
	ALLOW_MULTIPLE_ONLINE                  //3-允许多端登录，不互踢
	ONLY_ONE_WITHIN_TERMTYPES              //4-指定端互踢
	ONLY_ONE_WITHOUT_TERMTYPES             //5-指定端不互踢
	ONLY_ONE_WITHIN_TERMTYPE_GROUPS        //6-指定端互踢多分组
)
