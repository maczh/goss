package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const COLLECTION_WECHAT_USER_INFO = "WechatUserInfo"

func InsertWechatUserInfo(wechatUserInfo model.WechatUserInfo) (model.WechatUserInfo, error) {
	wechatUserInfo.ID = bson.NewObjectId()
	wechatUserInfo.BindTime = utils.ToDateTimeString(time.Now())
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.WechatUserInfo{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_WECHAT_USER_INFO).Insert(&wechatUserInfo)
	if err != nil {
		return model.WechatUserInfo{}, err
	}
	return wechatUserInfo, nil
}

func GetWechatUserInfo(userId, unionId string) (model.WechatUserInfo, error) {
	var wechatUserInfo model.WechatUserInfo
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.WechatUserInfo{}, errors.New("MongoDB连接异常")
	}
	query := bson.M{}
	if userId != "" {
		query["userId"] = userId
	}
	if unionId != "" {
		query["unionId"] = unionId
	}
	err := mongo.C(COLLECTION_WECHAT_USER_INFO).Find(&query).One(&wechatUserInfo)
	if err != nil && err.Error() != "not found" {
		return model.WechatUserInfo{}, err
	}
	return wechatUserInfo, nil
}

func UpdateWechatUserInfo(wechatUserInfo model.WechatUserInfo) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_WECHAT_USER_INFO).UpdateId(wechatUserInfo.ID, &wechatUserInfo)
	if err != nil {
		return err
	}
	return nil
}

func DeleteWechatUserInfo(wechatUserInfo model.WechatUserInfo) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_WECHAT_USER_INFO).RemoveId(wechatUserInfo.ID)
	if err != nil {
		return err
	}
	return nil
}
