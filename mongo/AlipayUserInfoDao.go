package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_ALIPAY_USER_INFO = "AlipayUserInfo"

func InsertAlipayUserInfo(alipayUserInfo model.AlipayUserInfo) (model.AlipayUserInfo, error) {
	alipayUserInfo.ID = bson.NewObjectId()
	alipayUserInfo.BindTime = utils.ToDateTimeString(time.Now())
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.AlipayUserInfo{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_ALIPAY_USER_INFO).Insert(&alipayUserInfo)
	if err != nil {
		return model.AlipayUserInfo{}, err
	}
	return alipayUserInfo, nil
}

func GetAlipayUserInfo(userId, alipayId string) (model.AlipayUserInfo, error) {
	var alipayUserInfo model.AlipayUserInfo
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.AlipayUserInfo{}, errors.New("MongoDB连接异常")
	}
	query := bson.M{}
	if userId != "" {
		query["userId"] = userId
	}
	if alipayId != "" {
		query["alipayId"] = alipayId
	}
	err := mongo.C(COLLECTION_ALIPAY_USER_INFO).Find(&query).One(&alipayUserInfo)
	if err != nil && err.Error() != "not found" {
		return model.AlipayUserInfo{}, err
	}
	return alipayUserInfo, nil
}

func UpdateAlipayUserInfo(alipayUserInfo model.AlipayUserInfo) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_ALIPAY_USER_INFO).UpdateId(alipayUserInfo.ID, &alipayUserInfo)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAlipayUserInfo(alipayUserInfo model.AlipayUserInfo) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_ALIPAY_USER_INFO).RemoveId(alipayUserInfo.ID)
	if err != nil {
		return err
	}
	return nil
}
