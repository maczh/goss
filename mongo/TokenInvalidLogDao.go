package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_INVALID_LOG = "TokenInvalidLog"

func InsertInvalidLog(invalidLog model.TokenInvalidLog) (model.TokenInvalidLog, error) {
	invalidLog.ID = bson.NewObjectId()
	invalidLog.Time = utils.ToDateTimeString(time.Now())
	invalidLog.Token = invalidLog.TokenInfo.Token
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.TokenInvalidLog{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_INVALID_LOG).Insert(&invalidLog)
	if err != nil {
		return model.TokenInvalidLog{}, err
	}
	return invalidLog, nil
}

func GetTokenInvalidLog(token string) (model.TokenInvalidLog, error) {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.TokenInvalidLog{}, errors.New("MongoDB连接异常")
	}
	var invalidLog model.TokenInvalidLog
	err := mongo.C(COLLECTION_INVALID_LOG).Find(&bson.M{"token": token}).One(&invalidLog)
	if err != nil && err.Error() != "not found" {
		return model.TokenInvalidLog{}, err
	}
	return invalidLog, nil
}
