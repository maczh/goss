package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_USER_DEVICE = "UserDevice"

func InsertUserDevice(userDevice model.UserDevice) (model.UserDevice, error) {
	userDevice.ID = bson.NewObjectId()
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserDevice{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_DEVICE).Insert(&userDevice)
	if err != nil {
		return model.UserDevice{}, err
	}
	return userDevice, nil
}

func GetUserDevice(userId, deviceId string) (model.UserDevice, error) {
	var userDevice model.UserDevice
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserDevice{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_DEVICE).Find(&bson.M{"userId": userId, "deviceId": deviceId}).One(&userDevice)
	if err != nil && err.Error() != "not found" {
		return model.UserDevice{}, err
	}
	return userDevice, nil
}

func UpdateUserDevice(userDevice model.UserDevice) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_DEVICE).UpdateId(userDevice.ID, &userDevice)
	if err != nil {
		return err
	}
	return nil
}

func ListUserDevice(userId string) ([]model.UserDevice, error) {
	userDevices := make([]model.UserDevice, 0)
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return nil, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_DEVICE).Find(&bson.M{"userId": userId}).All(&userDevices)
	if err != nil {
		return nil, err
	}
	return userDevices, nil
}
