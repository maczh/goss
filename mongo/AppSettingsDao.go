package mongo

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/mgcache"
	"github.com/maczh/mgconfig"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_APP_SETTINGS = "AppSettings"

func InsertAppSettings(appSettings model.AppSettings) (model.AppSettings, error) {
	appSettings.ID = bson.NewObjectId()
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.AppSettings{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_APP_SETTINGS).Insert(&appSettings)
	if err != nil {
		return model.AppSettings{}, err
	}
	return appSettings, nil
}

func GetAppSettings(appId string) (model.AppSettings, error) {
	var appSettings model.AppSettings
	if mgcache.OnGetCache("goss").IsExist("appSettings:" + appId) {
		appSettingsCache, _ := mgcache.OnGetCache("goss").Value("appSettings:" + appId)
		appSettings = appSettingsCache.(model.AppSettings)
		return appSettings, nil
	}
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return model.AppSettings{}, errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_APP_SETTINGS).Find(&bson.M{"appId": appId}).One(&appSettings)
	if err != nil && err.Error() != "not found" {
		return model.AppSettings{}, err
	}
	mgcache.OnGetCache("goss").Add("appSettings:"+appId, appSettings, 0)
	return appSettings, nil
}

func UpdateAppSettings(appSettings model.AppSettings) error {
	mongo := mgconfig.GetMongoConnection()
	defer mgconfig.ReturnMongoConnection(mongo)
	if mongo == nil {
		return errors.New("MongoDB连接异常")
	}
	err := mongo.C(COLLECTION_APP_SETTINGS).UpdateId(appSettings.ID, &appSettings)
	if err != nil {
		return err
	}
	mgcache.OnGetCache("goss").Delete("appSettings:" + appSettings.AppId)
	return nil
}

func UpsertAppSettings(appSettings model.AppSettings) (model.AppSettings, error) {
	appSettings, err := GetAppSettings(appSettings.AppId)
	//if err != nil {
	//	return appSettings, err
	//}
	if appSettings.AppId == appSettings.AppId {
		err = UpdateAppSettings(appSettings)
		return appSettings, err
	} else {
		return InsertAppSettings(appSettings)
	}
}
