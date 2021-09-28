package model

type AppInfo struct {
	Id       int    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	AppId    string `json:"app_id" gorm:"column:app_id;unique"`
	AppKey   string `json:"app_key" gorm:"column:app_key"`
	AppName  string `json:"app_name" gorm:"column:app_name"`
	Descript string `json:"descript" gorm:"column:descript"`
	Status   int    `json:"status" gorm:"column:status"`
}
