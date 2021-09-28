package mysql

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
)

const TABLE_APP_INFO = "app_info"

func InsertAppInfo(app model.AppInfo) (model.AppInfo, error) {
	if app.AppId == "" {
		app.AppId = utils.GetUUIDString()
	}
	if app.AppKey == "" {
		app.AppKey = utils.GetRandomCaseString(32)
	}
	app.Status = 1
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		logs.Error("MySQL数据库连接异常")
		return model.AppInfo{}, errors.New("MySQL数据库连接异常")
	}
	err := db.Table(TABLE_APP_INFO).Create(&app).Error
	if err != nil {
		return model.AppInfo{}, err
	}
	return app, nil
}

func GetAppInfoByAppId(appId string) (model.AppInfo, error) {
	var app model.AppInfo
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return model.AppInfo{}, errors.New("MySQL数据库连接异常")
	}
	db.Table(TABLE_APP_INFO).Where("app_id = ?", appId).First(&app)
	return app, nil
}

func UpdateAppInfo(app model.AppInfo) error {
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return errors.New("MySQL数据库连接异常")
	}
	err := db.Table(TABLE_APP_INFO).Where("id = ?", app.Id).Update(&app).Error
	return err
}

func DeleteAppInfo(id int) error {
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return errors.New("MySQL数据库连接异常")
	}
	err := db.Table(TABLE_APP_INFO).Delete(nil, "id = ?", id).Error
	return err
}
