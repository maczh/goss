package redis

import (
	"fmt"
	"github.com/maczh/mgconfig"
)

const KEY_USER_ID = "user:id:current"

func GenerateNewUserId() string {
	goredis := mgconfig.GetRedisConnection()
	defer mgconfig.ReturnRedisConnection(goredis)
	if goredis == nil {
		return ""
	}
	return fmt.Sprintf("%d", goredis.Incr(KEY_USER_ID).Val())
}
