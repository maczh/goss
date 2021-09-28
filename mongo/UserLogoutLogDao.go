package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_LOGOUT_LOG = "UserLogoutLog"

func InsertLogoutLog(logoutLog model.UserLogoutLog) (model.UserLogoutLog, error) {
	logoutLog.ID = bson.NewObjectId()
	logoutLog.Time = utils.ToDateTimeString(time.Now())
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserLogoutLog{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_LOGOUT_LOG).Insert(&logoutLog)
	if err != nil {
		return model.UserLogoutLog{}, err
	}
	return logoutLog, nil
}
