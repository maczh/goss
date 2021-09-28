package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_KICK_LOG = "UserKickLog"

func InsertKickLog(kickLog model.UserKickLog) (model.UserKickLog, error) {
	kickLog.ID = bson.NewObjectId()
	kickLog.Time = utils.ToDateTimeString(time.Now())
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserKickLog{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_KICK_LOG).Insert(&kickLog)
	if err != nil {
		return model.UserKickLog{}, err
	}
	return kickLog, nil
}
