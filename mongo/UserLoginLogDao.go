package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_LOGIN_LOG = "UserLoginLog"

func InsertLoginLog(loginLog model.UserLoginLog) (model.UserLoginLog, error) {
	loginLog.ID = bson.NewObjectId()
	loginLog.Time = utils.ToDateTimeString(time.Now())
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserLoginLog{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_LOGIN_LOG).Insert(&loginLog)
	if err != nil {
		return model.UserLoginLog{}, err
	}
	return loginLog, nil
}
