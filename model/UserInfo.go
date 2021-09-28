package model

import (
	"github.com/maczh/utils"
)

type UserInfo struct {
	Id         int    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UserId     string `json:"user_id" gorm:"column:user_id;unique"`
	UserName   string `json:"user_name" gorm:"column:user_name"`
	Password   string `json:"password" gorm:"column:password"`
	Mobile     string `json:"mobile" gorm:"column:mobile;unique"`
	Email      string `json:"email" gorm:"column:email"`
	Image      string `json:"image" gorm:"column:image"`
	CreateTime string `json:"create_time" gorm:"column:create_time"`
	Descript   string `json:"descript" gorm:"column:descript"`
	RealName   string `json:"real_name" gorm:"column:real_name"`
	IdCardNo   string `json:"id_card_no" gorm:"column:id_card_no"`
	Status     int    `json:"status" gorm:"column:status"`
}

func (u *UserInfo) GetPassword() string {
	pwd, _ := utils.AESBase64Decrypt(u.Password, utils.MD5Encode(u.UserId+u.Mobile), []byte{23, 67, 172, 65, 37, 210, 03, 82, 95, 183, 192, 236, 17, 72, 63, 166})
	return pwd
}

func (u *UserInfo) SetPassword(pwd string) {
	u.Password, _ = utils.AESBase64Encrypt(pwd, utils.MD5Encode(u.UserId+u.Mobile), []byte{23, 67, 172, 65, 37, 210, 03, 82, 95, 183, 192, 236, 17, 72, 63, 166})
}

func (u *UserInfo) VerifyPassword(pwd string) bool {
	return pwd == u.GetPassword()
}
