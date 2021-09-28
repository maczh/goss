package expire

import (
	"github.com/maczh/goss/constant"
	"github.com/maczh/goss/logic"
	"github.com/maczh/goss/redis"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"strconv"
	"strings"
)

func TokenExpiredListener() {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		logs.Error("Redis connection failed")
		return
	}
	logs.Debug("侦听过期消息: __keyevent@{}__:expired", goredis.Options().DB)
	pubsub := goredis.PSubscribe("__keyevent@" + strconv.Itoa(goredis.Options().DB) + "__:expired")
	_, err := pubsub.Receive()
	if err != nil {
		logs.Error("订阅错误:" + err.Error())
		return
	}
	ch := pubsub.Channel()
	for msg := range ch {
		if !strings.Contains(msg.Payload, redis.KEY_USER_TOKEN_EXPIRE) {
			continue
		}
		token := strings.ReplaceAll(msg.Payload, redis.KEY_USER_TOKEN_EXPIRE, "")
		logs.Debug("发现token={}过期，正在进行过期处理", token)
		_ = logic.DeleteUserToken(token, "过期失效", constant.IT_EXPIRED)
	}
}
