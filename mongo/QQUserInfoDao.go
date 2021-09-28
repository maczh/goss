package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_QQ_USER_INFO = "QQUserInfo"

func InsertQQUserInfo(qqUserInfo model.QQUserInfo) (model.QQUserInfo, error) {
	qqUserInfo.ID = bson.NewObjectId()
	qqUserInfo.BindTime = utils.ToDateTimeString(time.Now())
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.QQUserInfo{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_QQ_USER_INFO).Insert(&qqUserInfo)
	if err != nil {
		return model.QQUserInfo{}, err
	}
	return qqUserInfo, nil
}

func GetQQUserInfo(userId, openId string) (model.QQUserInfo, error) {
	var qqUserInfo model.QQUserInfo
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.QQUserInfo{}, errors.New("MongoDB连接异常")
	}
	query := bson.M{}
	if userId != "" {
		query["userId"] = userId
	}
	if openId != "" {
		query["openId"] = openId
	}
	err := mongo.C(COLLECTION_QQ_USER_INFO).Find(&query).One(&qqUserInfo)
	if err != nil && err.Error() != "not found" {
		return model.QQUserInfo{}, err
	}
	return qqUserInfo, nil
}

func UpdateQQUserInfo(qqUserInfo model.QQUserInfo) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_QQ_USER_INFO).UpdateId(qqUserInfo.ID, &qqUserInfo)
	if err != nil {
		return err
	}
	return nil
}

func DeleteQQUserInfo(qqUserInfo model.QQUserInfo) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_QQ_USER_INFO).RemoveId(qqUserInfo.ID)
	if err != nil {
		return err
	}
	return nil
}
