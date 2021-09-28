package mysql

import (
	"errors"
	"github.com/maczh/goss/model"
	"github.com/maczh/goss/redis"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/utils"
	"time"
)

const TABLE_USER_INFO = "user_info"

func InsertUserInfo(user model.UserInfo) (model.UserInfo, error) {
	if user.Mobile == "" {
		return user, errors.New("用户手机号码不可为空")
	}
	userInfo, err := GetUserInfoByMobile(user.Mobile)
	if err != nil {
		return user, err
	}
	if userInfo.Mobile != "" {
		return user, errors.New("手机号对应的账号已存在，不允许重复注册")
	}
	if user.UserId == "" {
		user.UserId = redis.GenerateNewUserId()
	}
	if user.Password == "" {
		user.SetPassword(utils.GetRandomCaseString(12))
	} else {
		userInfo.SetPassword(user.Password)
	}
	user.CreateTime = utils.ToDateTimeString(time.Now())
	user.Status = 1
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		logs.Error("MySQL数据库连接异常")
		return user, errors.New("MySQL数据库连接异常")
	}
	err = db.Table(TABLE_USER_INFO).Create(&user).Error
	if err != nil {
		return user, errors.New("账号数据入库错误:" + err.Error())
	}
	db.Table(TABLE_USER_INFO).Where("mobile = ?", user.Mobile).First(&user)
	return user, nil
}

func GetUserInfoByMobile(mobile string) (model.UserInfo, error) {
	var user model.UserInfo
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return model.UserInfo{}, errors.New("MySQL数据库连接异常")
	}
	db.Table(TABLE_USER_INFO).Where("mobile = ?", mobile).First(&user)
	return user, nil
}

func GetUserInfoByUserId(userId string) (model.UserInfo, error) {
	var user model.UserInfo
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return model.UserInfo{}, errors.New("MySQL数据库连接异常")
	}
	db.Table(TABLE_USER_INFO).Where("user_id = ?", userId).First(&user)
	return user, nil
}

func UpdateUserInfo(user model.UserInfo) error {
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return errors.New("MySQL数据库连接异常")
	}
	err := db.Table(TABLE_USER_INFO).Where("id = ?", user.Id).Update(&user).Error
	return err
}

func DeleteUserInfo(id int) error {
	db := mgconfig.GetMysqlConnection()
	defer mgconfig.ReturnMysqlConnection(db)
	if db == nil {
		return errors.New("MySQL数据库连接异常")
	}
	err := db.Table(TABLE_USER_INFO).Delete(nil, "id = ?", id).Error
	return err
}
