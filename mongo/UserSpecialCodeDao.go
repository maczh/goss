package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgconfig"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_USER_SPECIAL_CODE = "UserSpecialCode"

func InsertUserSpecialCode(userSpecialCode model.UserSpecialCode) (model.UserSpecialCode, error) {
	userSpecialCode.ID = bson.NewObjectId()
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserSpecialCode{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_SPECIAL_CODE).Insert(&userSpecialCode)
	if err != nil {
		return model.UserSpecialCode{}, err
	}
	return userSpecialCode, nil
}

func GetUserSpecialCode(userId string) (model.UserSpecialCode, error) {
	var userSpecialCode model.UserSpecialCode
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.UserSpecialCode{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_SPECIAL_CODE).Find(&bson.M{"userId": userId}).One(&userSpecialCode)
	if err != nil && err.Error() != "not found" {
		return model.UserSpecialCode{}, err
	}
	return userSpecialCode, nil
}

func UpdateUserSpecialCode(userSpecialCode model.UserSpecialCode) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_USER_SPECIAL_CODE).UpdateId(userSpecialCode.ID, &userSpecialCode)
	if err != nil {
		return err
	}
	return nil
}

func UpsertUserSpecialCode(userSpecialCode model.UserSpecialCode) (model.UserSpecialCode, error) {
	specialCode, err := GetUserSpecialCode(userSpecialCode.UserId)
	if err != nil && err.Error() != "not found" {
		return userSpecialCode, err
	}
	if specialCode.UserId == userSpecialCode.UserId {
		err = UpdateUserSpecialCode(userSpecialCode)
		return userSpecialCode, err
	} else {
		return InsertUserSpecialCode(userSpecialCode)
	}
}
